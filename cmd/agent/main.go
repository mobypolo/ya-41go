package main

import (
	"fmt"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/agent/sources"
	"io"
	"log"
	"net/http"
	"time"
)

import _ "github.com/mobypolo/ya-41go/internal/agent/sources"

var (
	serverAddress  = "http://localhost:8080"
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	metricsChan := make(chan []agent.Metric, 1)

	go func() {
		for {
			metrics, err := agent.CollectAll()
			if err == nil {
				metricsChan <- metrics
			}
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		ticker := time.NewTicker(reportInterval)
		defer ticker.Stop()

		for range ticker.C {
			var metrics []agent.Metric
			select {
			case metrics = <-metricsChan:
				sources.PollCount++
			default:
				metrics = nil
			}

			for _, m := range metrics {
				sendMetric(m)
			}
		}
	}()

	log.Println("Agent started")
	select {}
}

func sendMetric(m agent.Metric) {
	url := fmt.Sprintf("%s/update/%s/%s/%v", serverAddress, m.Type, m.Name, m.Value)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Println("build request error:", err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("request error:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("close body error:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("server responded with %s for %s", resp.Status, url)
	}
}
