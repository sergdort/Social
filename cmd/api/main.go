package main

import (
	"github.com/sergdort/Social/internal/db"
	"github.com/sergdort/Social/internal/env"
	"github.com/sergdort/Social/internal/store"
	"log"
)

const version = "0.0.0.1"

func main() {
	var cfg = config{
		address: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	var database, err = db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer database.Close()

	log.Println("database connection established")

	var s = store.NewStorage(database)

	var app = &application{
		config: cfg,
		store:  s,
	}

	log.Fatal(app.run(app.mount()))
}
