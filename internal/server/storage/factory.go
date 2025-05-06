package storage

import (
	"github.com/mobypolo/ya-41go/cmd"
	"github.com/mobypolo/ya-41go/internal/server/db"
	"time"
)

func MakeStorage(cfg cmd.Config) *PersistentStorage {
	if cfg.DatabaseDSN != "" {
		db.InitPostgres(cfg.DatabaseDSN)
		return NewPersistentStorageWithPostgres()
	}

	ps := NewPersistentStorage(
		cfg.FileStoragePath,
		time.Duration(cfg.StoreInterval)*time.Second,
		cfg.RestoreOnStart,
	)
	return ps
}
