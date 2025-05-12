package main

import (
	"context"
	"expvar"
	"github.com/redis/go-redis/v9"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/business/platform/jwt"
	"github.com/sergdort/Social/business/platform/mailer"
	"github.com/sergdort/Social/business/platform/store"
	"github.com/sergdort/Social/business/platform/store/cache"
	"github.com/sergdort/Social/cmd/api/debug"
	"github.com/sergdort/Social/foundation/logger"
	"github.com/sergdort/Social/foundation/otel"
	"github.com/sergdort/Social/internal/auth"
	"github.com/sergdort/Social/internal/db"
	"github.com/sergdort/Social/internal/env"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const version = "0.0.0.1"

//	@title			Go Social
//	@description	API for Go Social.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {
	var cfg = config{
		address:         env.GetString("ADDR", ":8080"),
		debugHost:       env.GetString("DEBUG_HOST", ":8090"),
		shutDownTimeout: time.Duration(env.GetInt("SHUTDOWN_TIMEOUT", 20)),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", true),
		},
		env:    env.GetString("ENV", "development"),
		apiURL: env.GetString("EXTERNAL_URL", "http://localhost:8080"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("SENDGRID_FROM_EMAIL", "<EMAIL>"),
			sendGridConfig: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		frontEndURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		auth: authConfig{
			basic: basicAuthConfig{
				username: env.GetString("AUTH_BASIC_USERNAME", "admin"),
				password: env.GetString("AUTH_BASIC_PASSWORD", "admin"),
			},
			jwt: jwtAuthConfig{
				secret:    env.GetString("JWT_SECRET", "secret"),
				exp:       time.Hour * 24 * 7,
				tokenHost: env.GetString("JWT_TOKEN_HOST", "social"),
			},
		},
		serviceName: env.GetString("SERVICE_NAME", "social"),
	}
	ctx := context.Background()
	var log *logger.Logger
	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}
	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}
	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SOCIAL", traceIDFn, events)
	// Mailer
	mail := mailer.NewSendgridMailer(cfg.mail.fromEmail, cfg.mail.sendGridConfig.apiKey)

	// Database
	var database, err = db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}

	defer database.Close()

	log.Info(
		ctx,
		"DB connected",
		"addr", cfg.db.addr,
		"maxOpenConns", cfg.db.maxOpenConns,
		"maxIdleConns", cfg.db.maxIdleConns,
		"maxIdleTime", cfg.db.maxIdleTime,
	)
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		log.Info(ctx, "Redis connected", "addr", cfg.redisCfg.addr)
	}
	cacheStorage := cache.NewStorage(rdb)

	s := store.NewStorage(database)

	var authenticator = auth.NewJWTAutheticator(
		cfg.auth.jwt.secret,
		cfg.auth.jwt.tokenHost,
		cfg.auth.jwt.tokenHost,
	)

	jwtAuth := jwt.NewJWTAutheticator(
		cfg.auth.jwt.secret,
		cfg.auth.jwt.tokenHost,
		cfg.auth.jwt.tokenHost,
		cfg.auth.jwt.tokenHost,
		cfg.auth.jwt.exp,
	)

	var app = &application{
		config:        cfg,
		store:         s,
		logger:        log,
		mailer:        mail,
		authenticator: authenticator,
		cache:         cacheStorage,
		useCase: useCases{
			Users: domain.NewUsersUseCase(cacheStorage.Users, s.Users),
			Auth: domain.NewAuthUseCase(
				domain.AuthConfig{
					InvitationExp: cfg.mail.exp,
					FrontendURL:   cfg.frontEndURL,
				},
				s.Roles,
				s.Users,
				jwtAuth,
			),
		},
	}
	// TODO: Pass build type
	expvar.NewString("build").Set("develop")

	go func() {
		log.Info(
			ctx,
			"debug v1 router started",
			"host", cfg.debugHost,
		)
		if err := http.ListenAndServe(cfg.debugHost, debug.Mux()); err != nil {
			log.Error(ctx, "debug v1 router stopped", "host", cfg.debugHost, "err", err)
		}
	}()

	server := app.makeServer(app.mount(ctx, log))
	serverErrors := make(chan error, 1)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info(ctx, "api router started", "host", cfg.apiURL)

		serverErrors <- server.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		log.Info(ctx, "api router stopped", "host", cfg.apiURL, "err", err)
	case sig := <-shutdown:
		log.Info(ctx, "shutdown started", "status", sig)
		defer log.Info(ctx, "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.shutDownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			log.Error(ctx, "could not stop server gracefully: %w", err)
		}
	}
}
