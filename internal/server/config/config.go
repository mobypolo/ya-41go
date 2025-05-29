package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/mobypolo/ya-41go/internal/shared/logger"
	"github.com/spf13/pflag"
	"log"
)

type Config struct {
	Address         string            `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval   int               `env:"STORE_INTERVAL" envDefault:"300"`
	FileStoragePath string            `env:"FILE_STORAGE_PATH" envDefault:"metrics.json"`
	RestoreOnStart  bool              `env:"RESTORE" envDefault:"true"`
	DatabaseDSN     string            `env:"DATABASE_DSN" envDefault:""`
	Key             string            `env:"KEY" envDefault:""`
	LogMode         logger.ModeLogger `env:"LOG_MODE" envDefault:"dev"`
}

func ParseFlags() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("env parse error: %v", err)
	}

	pflag.StringVarP(&cfg.Address, "address", "a", cfg.Address, "HTTP server address")
	pflag.IntVarP(&cfg.StoreInterval, "interval", "i", cfg.StoreInterval, "Store interval (seconds)")
	pflag.StringVarP(&cfg.FileStoragePath, "file", "f", cfg.FileStoragePath, "File Storage Path")
	pflag.BoolVarP(&cfg.RestoreOnStart, "restore", "r", cfg.RestoreOnStart, "Restore on load")
	pflag.StringVarP(&cfg.DatabaseDSN, "dsn", "d", cfg.DatabaseDSN, "PostgresSQL DSN")
	pflag.StringVarP(&cfg.Key, "key", "k", cfg.Key, "Secret key for HMAC SHA256")
	pflag.Parse()

	return cfg
}
