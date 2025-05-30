package main

import (
	"context"
	"github.com/mobypolo/ya-41go/internal/agent/bootstrap"
	"github.com/mobypolo/ya-41go/internal/agent/config"
	"log"
)

func main() {
	cfg := config.ParseFlags()
	ctx := context.Background()
	if err := bootstrap.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
