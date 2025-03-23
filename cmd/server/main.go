package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mobypolo/ya-41go/internal/server/route"
	"log"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/server/handler"

func main() {
	r := chi.NewRouter()
	route.MountInto(r)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
