package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/business/platform/mailer"
	"github.com/sergdort/Social/business/platform/store"
	cache2 "github.com/sergdort/Social/business/platform/store/cache"
	"github.com/sergdort/Social/cmd/api/debug"
	"github.com/sergdort/Social/internal/auth"
	"github.com/sergdort/Social/internal/db"
	"github.com/sergdort/Social/internal/env"
	"go.uber.org/zap"
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
	}
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

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
		logger.Fatal(err)
	}

	defer database.Close()

	logger.Infow(
		"DB connected",
		"addr", cfg.db.addr,
		"maxOpenConns", cfg.db.maxOpenConns,
		"maxIdleConns", cfg.db.maxIdleConns,
		"maxIdleTime", cfg.db.maxIdleTime,
	)
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache2.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Infow("Redis connected")
	}
	cacheStorage := cache2.NewStorage(rdb)

	s := store.NewStorage(database)

	var authenticator = auth.NewJWTAutheticator(
		cfg.auth.jwt.secret,
		cfg.auth.jwt.tokenHost,
		cfg.auth.jwt.tokenHost,
	)

	var app = &application{
		config:        cfg,
		store:         s,
		logger:        logger,
		mailer:        mail,
		authenticator: authenticator,
		cache:         cacheStorage,
		useCase: useCases{
			Users: domain.NewUsersUseCase(cacheStorage.Users, s.Users),
		},
	}
	ctx := context.Background()
	// TODO: Pass build type
	expvar.NewString("build").Set("develop")

	go func() {
		logger.Infow(
			"debug v1 router started",
			"host", cfg.debugHost,
		)
		if err := http.ListenAndServe(cfg.debugHost, debug.Mux()); err != nil {
			logger.Errorw("debug v1 router stopped", "host", cfg.debugHost, "err", err)
		}
	}()

	serverErrors := make(chan error, 1)
	server := app.makeServer(app.mount())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infow("api router started", "host", cfg.apiURL)

		serverErrors <- server.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		logger.Infow("api router stopped", "host", cfg.apiURL, "err", err)
	case sig := <-shutdown:
		logger.Infow("shutdown started", "status", sig)
		defer logger.Infow("shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.shutDownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
}
