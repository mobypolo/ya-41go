package metrics

import (
	"github.com/mobypolo/ya-41go/internal/server/storage"
)

func init() {
	storage.RegisterProcessor(string(storage.GaugeType), NewGaugeProcessor())
	storage.RegisterProcessor(string(storage.CounterType), NewCounterProcessor())
}
