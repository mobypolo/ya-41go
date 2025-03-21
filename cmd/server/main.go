package main

import (
	"github.com/mobypolo/ya-41go/internal/router"
	"log"
	"net/http"
)

func main() {
	r := router.NewRouter()

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
