package agent

import "github.com/mobypolo/ya-41go/internal/server/storage"

type MetricSource interface {
	Name() string
	Type() storage.MetricType
	Collect() (interface{}, error)
}
