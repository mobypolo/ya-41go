package workers

import (
	"context"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/agent/config"
	"time"
)

func StartCollectors(_ context.Context, metricsChan chan []agent.Metric, cfg config.Config) {
	ticker := time.NewTicker(cfg.Poll)
	defer ticker.Stop()
	for range ticker.C {
		metrics, err := agent.CollectAll()
		if err == nil {
			metricsChan <- metrics
		}
	}
}
