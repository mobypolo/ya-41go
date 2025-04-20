package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mobypolo/ya-41go/cmd"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/agent/sources"
	"github.com/mobypolo/ya-41go/internal/shared/dto"
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
				//sendMetric(m)
				sendMetricJSON(m)
			}
		}
	}()

	log.Println("Agent started")
	select {}
}

func sendMetric(m agent.Metric) {
	serverAddress := fmt.Sprintf("http://%s", cmd.ServerAddress)

	url := fmt.Sprintf("%s/update_plain/%s/%s/%v", serverAddress, m.Type, m.Name, m.Value)
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

func sendMetricJSON(m agent.Metric) {
	serverAddress := fmt.Sprintf("http://%s", cmd.ServerAddress)
	url := fmt.Sprintf("%s/update/", serverAddress)

	var payload dto.Metrics
	payload.ID = m.Name
	payload.MType = string(m.Type)

	switch m.Type {
	case "gauge":
		if val, ok := m.Value.(float64); ok {
			payload.Value = &val
		} else {
			log.Printf("invalid gauge value: %v", m.Value)
			return
		}
	case "counter":
		switch v := m.Value.(type) {
		case int64:
			payload.Delta = &v
		case float64:
			val := int64(v)
			payload.Delta = &val
		default:
			log.Printf("invalid counter value type: %T", v)
			return
		}
	default:
		log.Printf("unknown metric type: %s", m.Type)
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal JSON: %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("build request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

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
