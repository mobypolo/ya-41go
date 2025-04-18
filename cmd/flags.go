package cmd

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type config struct {
	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
}

var (
	ServerAddress  string
	ReportInterval time.Duration
	PollInterval   time.Duration
)

func ParseFlags(app string) {
	var cfg config

	if err := env.Parse(&cfg); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to parse env: %v\n", err)
		if err != nil {
			log.Println(err)
		}
		os.Exit(1)
	}

	switch app {
	case "server":
		pflag.StringVarP(&ServerAddress, "address", "a", cfg.Address, "HTTP server address")

	case "agent":
		pflag.StringVarP(&ServerAddress, "address", "a", cfg.Address, "HTTP server address")
		report := pflag.IntP("report-interval", "r", cfg.ReportInterval, "Report interval (seconds)")
		poll := pflag.IntP("poll-interval", "p", cfg.PollInterval, "Poll interval (seconds)")

		pflag.Parse()

		ReportInterval = time.Duration(*report) * time.Second
		PollInterval = time.Duration(*poll) * time.Second

	default:
		_, err := fmt.Fprintf(os.Stderr, "Unknown app type: %s\n", app)
		if err != nil {
			log.Println(err)
		}
		os.Exit(1)
	}

	pflag.Parse()

	if len(pflag.Args()) > 0 {
		_, err := fmt.Fprintf(os.Stderr, "Unknown flags: %v\n", pflag.Args())
		if err != nil {
			log.Println(err)
		}
		os.Exit(1)
	}
}
