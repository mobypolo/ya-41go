package storage

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mobypolo/ya-41go/internal/server/config"
	"time"
)

func MakeStorage(cfg config.Config, db *pgxpool.Pool) *PersistentStorage {
	if cfg.DatabaseDSN != "" && db != nil {
		return NewPersistentStorageWithPostgres(db)
	}

	ps := NewPersistentStorage(
		cfg.FileStoragePath,
		time.Duration(cfg.StoreInterval)*time.Second,
		cfg.RestoreOnStart,
	)
	return ps
}
