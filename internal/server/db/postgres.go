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

func InitPostgres(dsn string) *pgxpool.Pool {
	if dsn == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolLocal, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.L().Error("failed to create db pool", zap.Error(err))
		return nil
	}

	if err := poolLocal.Ping(ctx); err != nil {
		logger.L().Error("failed to ping db", zap.Error(err))
		poolLocal.Close()
		return nil
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.L().Info("gorm connection to postgres failed")
		return nil
	}
	if err := gormDB.AutoMigrate(&Metric{}); err != nil {
		logger.L().Info("auto migration to postgres failed")
		return nil
	}

	logger.L().Info("connected to postgres")
	return poolLocal
}
