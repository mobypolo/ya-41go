package workers

import (
	"context"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/agent/config"
	"github.com/mobypolo/ya-41go/internal/agent/sender"
	"github.com/mobypolo/ya-41go/internal/shared/utils"
)

func StartSenderWorkers(ctx context.Context, metricTasks chan agent.Metric, cfg config.Config) {
	for i := 0; i < cfg.RateLimit; i++ {
		go func() {
			for metric := range metricTasks {
				_ = utils.RetryWithBackoff(ctx, 3, func() error {
					return sender.SendMetricJSON(metric, cfg)
				})
			}
		}()
	}
}
