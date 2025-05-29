package bootstrap

import (
	"context"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/agent/config"
	"github.com/mobypolo/ya-41go/internal/agent/sources"
	"github.com/mobypolo/ya-41go/internal/agent/workers"
	"log"
)

func Run(ctx context.Context, cfg config.Config) error {
	sources.StartRuntimeGaugeMetrics()

	metricsChan := make(chan []agent.Metric, 1)
	bufferSize := cfg.RateLimit * 10
	metricTasks := make(chan agent.Metric, bufferSize)

	go workers.StartCollectors(ctx, metricsChan, cfg)
	workers.StartSenderWorkers(ctx, metricTasks, cfg)
	go workers.StartBatchReporter(ctx, metricsChan, metricTasks, cfg)

	log.Println("Agent started")
	log.Println("Started agent with cfg : ", fmt.Sprintf("%+v", cfg))
	<-ctx.Done()
	return ctx.Err()
}
