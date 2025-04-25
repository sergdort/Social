package main

import (
	"github.com/sergdort/Social/business/platform/store"
	"github.com/sergdort/Social/internal/db"
	"github.com/sergdort/Social/internal/env"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	con, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	s := store.NewStorage(con)
	db.Seed(s, con)
}
