package workers

import (
	"context"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/agent/config"
	"github.com/mobypolo/ya-41go/internal/agent/sender"
	"github.com/mobypolo/ya-41go/internal/agent/sources"
	"github.com/mobypolo/ya-41go/internal/shared/utils"
	"time"
)

func StartBatchReporter(ctx context.Context, metricsChan chan []agent.Metric, metricTasks chan agent.Metric, cfg config.Config) {
	ticker := time.NewTicker(cfg.Report)
	defer ticker.Stop()

	for range ticker.C {
		var metrics []agent.Metric
		select {
		case metrics = <-metricsChan:
			sources.PollCount++
		default:
			continue
		}
		if len(metrics) == 0 {
			continue
		}

		for _, m := range metrics {
			metricTasks <- m
		}
		_ = utils.RetryWithBackoff(ctx, 3, func() error {
			return sender.SendMetricJSONBatch(metrics, cfg)
		})
	}

}
