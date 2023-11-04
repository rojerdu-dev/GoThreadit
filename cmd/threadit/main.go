package main

import (
	"github.com/rojerdu-dev/gothreadit/postgres"
	web "github.com/rojerdu-dev/gothreadit/web"
	"log"
	"net/http"
)

func main() {
	store, err := postgres.NewStore("postgres://postgres:pass@127.0.0.1:5432?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	h := web.NewHandler(store)
	http.ListenAndServe(":3000", h)
}
