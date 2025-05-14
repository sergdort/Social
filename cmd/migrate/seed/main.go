package main

import (
	db2 "github.com/sergdort/Social/business/platform/db"
	"github.com/sergdort/Social/business/platform/store"
	"github.com/sergdort/Social/foundation/env"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	con, err := db2.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	s := store.NewStorage(con)
	db2.Seed(s, con)
}
