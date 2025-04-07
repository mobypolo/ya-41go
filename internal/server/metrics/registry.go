package metrics

import (
	storage2 "github.com/mobypolo/ya-41go/internal/server/storage"
)

func init() {
	storage2.RegisterProcessor(string(storage2.GaugeType), NewGaugeProcessor())
	storage2.RegisterProcessor(string(storage2.CounterType), NewCounterProcessor())
}
