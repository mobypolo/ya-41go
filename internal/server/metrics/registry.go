package metrics

import (
	"github.com/mobypolo/ya-41go/internal/server/storage"
)

func init() {
	storage.RegisterProcessor(storage.GaugeType, NewGaugeProcessor())
	storage.RegisterProcessor(storage.CounterType, NewCounterProcessor())
}
