package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/mobypolo/ya-41go/internal/shared/logger"
	"github.com/spf13/pflag"
	"log"
	"time"
)

type Config struct {
	Address   string            `env:"ADDRESS" envDefault:"localhost:8080"`
	Key       string            `env:"KEY" envDefault:""`
	RateLimit int               `env:"RATE_LIMIT" envDefault:"3"`
	LogMode   logger.ModeLogger `env:"LOG_MODE" envDefault:"dev"`
	Report    time.Duration
	Poll      time.Duration
}

func ParseFlags() Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("env parse error: %v", err)
	}

	report := pflag.IntP("report-interval", "r", 10, "Report interval (seconds)")
	poll := pflag.IntP("poll-interval", "p", 2, "Poll interval (seconds)")
	pflag.StringVarP(&cfg.Address, "address", "a", cfg.Address, "HTTP server address")
	pflag.StringVarP(&cfg.Key, "key", "k", cfg.Key, "Secret key for HMAC SHA256")
	pflag.IntVarP(&cfg.RateLimit, "limit", "l", cfg.RateLimit, "Max concurrent outgoing requests")
	pflag.Parse()

	cfg.Report = time.Duration(*report) * time.Second
	cfg.Poll = time.Duration(*poll) * time.Second

	if cfg.RateLimit <= 0 {
		cfg.RateLimit = 3
		log.Println("Rate limit must be more than 0, fallback to 3")
	}

	return cfg
}
