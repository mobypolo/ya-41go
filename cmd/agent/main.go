package main

import (
	"fmt"
	"github.com/mobypolo/ya-41go/cmd"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/agent/sources"
	"log"
	"net/http"
	"time"
)

func main() {
	cmd.ParseFlags("agent")
	metricsChan := make(chan []agent.Metric, 1)

	go func() {
		for {
			metrics, err := agent.CollectAll()
			if err == nil {
				metricsChan <- metrics
			}
			time.Sleep(cmd.PollInterval)
		}
	}()

	go func() {
		ticker := time.NewTicker(cmd.ReportInterval)
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
	serverAddress := fmt.Sprintf("http://%s", cmd.ServerAddress)

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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("server responded with %s for %s", resp.Status, url)
	}
}
