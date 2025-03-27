package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mobypolo/ya-41go/cmd"
	"github.com/mobypolo/ya-41go/internal/server/route"
	"log"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/server/handler"

func main() {
	cmd.ParseFlags("server")
	r := chi.NewRouter()
	route.MountInto(r)

	log.Println(fmt.Sprintf("Server started on %s", cmd.ServerAddress))
	if err := http.ListenAndServe(cmd.ServerAddress, r); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
