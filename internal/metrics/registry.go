package metrics

import "github.com/mobypolo/ya-41go/internal/storage"

func init() {
	storage.RegisterProcessor(string(storage.GaugeType), gaugeProcessor{})
	storage.RegisterProcessor(string(storage.CounterType), counterProcessor{})
}
