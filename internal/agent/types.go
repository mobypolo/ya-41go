package agent

import "github.com/mobypolo/ya-41go/internal/server/storage"

type MetricType string

type Metric struct {
	Name  string
	Type  storage.MetricType
	Value interface{} // float64 или int64
}
