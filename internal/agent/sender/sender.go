package sender

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/agent/config"
	"github.com/mobypolo/ya-41go/internal/agent/helpers"
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"github.com/mobypolo/ya-41go/internal/shared/dto"
	"github.com/mobypolo/ya-41go/internal/shared/utils"
	"log"
	"net/http"
)

// sendMetric
func _(m agent.Metric, cfg config.Config) {
	serverAddress := fmt.Sprintf("http://%s", cfg.Address)

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

func SendMetricJSON(m agent.Metric, cfg config.Config) error {
	serverAddress := fmt.Sprintf("http://%s", cfg.Address)
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
			return fmt.Errorf("invalid gauge value: %v", m.Value)
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
			return fmt.Errorf("invalid counter value type: %T", v)
		}
	default:
		log.Printf("unknown metric type: %s", m.Type)
		return fmt.Errorf("unknown metric type: %s", m.Type)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal JSON: %v", err)
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	compressedBody, err := helpers.CompressRequest(body)
	if err != nil {
		log.Println("compression error:", err)
		return fmt.Errorf("compression error: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, compressedBody)
	if err != nil {
		log.Println("build request error:", err)
		return fmt.Errorf("build request error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("request error:", err)
		return fmt.Errorf("request error: %s", err)
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

func SendMetricJSONBatch(metrics []agent.Metric, cfg config.Config) error {
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

	serverAddress := fmt.Sprintf("http://%s", cfg.Address)
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
