package dto

import "github.com/mobypolo/ya-41go/internal/server/storage"

type Metrics struct {
	ID    string             `json:"id"`
	MType storage.MetricType `json:"type"`
	Delta *int64             `json:"delta,omitempty"`
	Value *float64           `json:"value,omitempty"`
}
