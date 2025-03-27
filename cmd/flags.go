package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"
)

var (
	ServerAddress  string
	ReportInterval time.Duration
	PollInterval   time.Duration
)

func ParseFlags(app string) {
	switch app {
	case "server":
		pflag.StringVarP(&ServerAddress, "address", "a", "localhost:8080", "HTTP server address")

	case "agent":
		pflag.StringVarP(&ServerAddress, "address", "a", "localhost:8080", "HTTP server address")
		report := pflag.IntP("report-interval", "r", 10, "Report interval (seconds)")
		poll := pflag.IntP("poll-interval", "p", 2, "Poll interval (seconds)")

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
