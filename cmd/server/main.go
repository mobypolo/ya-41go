package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mobypolo/ya-41go/cmd"
	"github.com/mobypolo/ya-41go/internal/server/db"
	"github.com/mobypolo/ya-41go/internal/server/middleware"
	"github.com/mobypolo/ya-41go/internal/server/route"
	"github.com/mobypolo/ya-41go/internal/server/service"
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"github.com/mobypolo/ya-41go/internal/shared/logger"
	"go.uber.org/zap"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/server/handler"

func main() {
	cfg := cmd.ParseFlags("server")
	logger.Init(cfg.ModeLogger)

	dbInstancePool := db.InitPostgres(cfg.DatabaseDSN)
	store := storage.MakeStorage(cfg, dbInstancePool)

	defer store.Stop()
	service.SetMetricService(service.NewMetricService(store))

	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.GzipDecompressMiddleware)
	r.Use(middleware.GzipCompressMiddleware)

	route.RegisterAllRoutes(dbInstancePool)
	route.MountInto(r)

	logger.L().Info("Server started", zap.String("addr", cmd.ServerAddress))
	if err := http.ListenAndServe(cmd.ServerAddress, r); err != nil {
		logger.L().Fatal("server error", zap.Error(err))
	}
}
