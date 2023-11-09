package main

import (
	"github.com/rojerdu-dev/gothreadit/postgres"
	web "github.com/rojerdu-dev/gothreadit/web"
	"log"
	"net/http"
)

func main() {
	var dsn = "postgres://postgres:pass@127.0.0.1:5432?sslmode=disable"

	store, err := postgres.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	sessions, err := web.NewSessionManager(dsn)
	if err != nil {
		log.Fatal(err)
	}

	csrfKey := []byte("01234567890123456789012345678901")
	h := web.NewHandler(store, sessions, csrfKey)
	http.ListenAndServe(":3000", h)
}
