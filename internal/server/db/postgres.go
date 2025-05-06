package db

import (
	"context"
	"github.com/mobypolo/ya-41go/internal/shared/logger"
	"go.uber.org/zap"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitPostgres(dsn string) {
	if dsn == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.L().Error("failed to create db pool", zap.Error(err))
		return
	}

	if err := pool.Ping(ctx); err != nil {
		logger.L().Error("failed to ping db", zap.Error(err))
		pool.Close()
		return
	}

	logger.L().Info("connected to postgres")
	Pool = pool
}
