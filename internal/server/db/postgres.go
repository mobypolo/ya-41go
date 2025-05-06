package db

import (
	"context"
	"github.com/mobypolo/ya-41go/internal/shared/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitPostgres(dsn string) {
	if dsn == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolLocal, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.L().Error("failed to create db pool", zap.Error(err))
		return
	}

	if err := poolLocal.Ping(ctx); err != nil {
		logger.L().Error("failed to ping db", zap.Error(err))
		poolLocal.Close()
		return
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.L().Info("gorm connection to postgres failed")
		return
	}
	if err := gormDB.AutoMigrate(&Metric{}); err != nil {
		logger.L().Info("auto migration to postgres failed")
		return
	}

	logger.L().Info("connected to postgres")
	pool = poolLocal
}

func Pool() *pgxpool.Pool {
	return pool
}
