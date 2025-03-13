package main

import (
	"github.com/sergdort/Social/internal/db"
	"github.com/sergdort/Social/internal/env"
	"github.com/sergdort/Social/internal/mailer"
	"github.com/sergdort/Social/internal/store"
	"go.uber.org/zap"
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
		address: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
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

	var s = store.NewStorage(database)

	var app = &application{
		config: cfg,
		store:  s,
		logger: logger,
		mailer: mail,
	}

	logger.Fatal(app.run(app.mount()))
}
