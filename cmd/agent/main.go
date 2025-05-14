package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/cmd"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/agent/helpers"
	"github.com/mobypolo/ya-41go/internal/agent/sources"
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"github.com/mobypolo/ya-41go/internal/shared/dto"
	"github.com/mobypolo/ya-41go/internal/shared/utils"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := cmd.ParseFlags("agent")
	log.Println("Started agent with cfg : ", fmt.Sprintf("%+v", cfg))
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

			//for _, m := range metrics {
			//sendMetric(m)
			//sendMetricJSON(m)
			//}
			_ = utils.RetryWithBackoff(context.Background(), 3, func() error {
				return sendMetricJSONBatch(metrics, cfg)
			})
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
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

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

// sendMetricJSON
func _(m agent.Metric) {
	serverAddress := fmt.Sprintf("http://%s", cmd.ServerAddress)
	url := fmt.Sprintf("%s/update/", serverAddress)

	var payload dto.Metrics
	payload.ID = m.Name
	payload.MType = m.Type

	switch m.Type {
	case storage.GaugeType:
		if val, ok := m.Value.(float64); ok {
			payload.Value = &val
		} else {
			log.Printf("invalid gauge value: %v", m.Value)
			return
		}
	case storage.CounterType:
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

	compressedBody, err := helpers.CompressRequest(body)
	if err != nil {
		log.Println("compression error:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, compressedBody)
	if err != nil {
		log.Println("build request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

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

func sendMetricJSONBatch(metrics []agent.Metric, cfg cmd.Config) error {
	if len(metrics) == 0 {
		return errors.New("empty batch")
	}

	var batch []dto.Metrics

	for _, m := range metrics {
		var payload dto.Metrics
		payload.ID = m.Name
		payload.MType = m.Type

		switch m.Type {
		case storage.GaugeType:
			if val, ok := m.Value.(float64); ok {
				payload.Value = &val
			} else {
				log.Printf("invalid gauge value: %v", m.Value)
				continue
			}
		case storage.CounterType:
			switch v := m.Value.(type) {
			case int64:
				payload.Delta = &v
			case float64:
				val := int64(v)
				payload.Delta = &val
			default:
				log.Printf("invalid counter value type: %T", v)
				continue
			}
		default:
			log.Printf("unknown metric type: %s", m.Type)
			continue
		}

		batch = append(batch, payload)
	}

	if len(batch) == 0 {
		return errors.New("empty batch")
	}

	body, err := json.Marshal(batch)
	if err != nil {
		log.Printf("failed to marshal JSON batch: %v", err)
		return err
	}

	compressedBody, err := helpers.CompressRequest(body)
	if err != nil {
		log.Println("compression error:", err)
		return err
	}

	serverAddress := fmt.Sprintf("http://%s", cmd.ServerAddress)
	url := fmt.Sprintf("%s/updates/", serverAddress)

	req, err := http.NewRequest(http.MethodPost, url, compressedBody)
	if err != nil {
		log.Println("build request error:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	if cfg.Key != "" {
		hash := utils.HashBody(body, cfg.Key)
		req.Header.Set("HashSHA256", hash)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("request error:", err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("server responded with %s for %s", resp.Status, url)
	}

	return nil
}
