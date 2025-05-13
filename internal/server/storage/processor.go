package storage

import (
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/server/repositories"
)

type MetricProcessor interface {
	ValidateName(name string) error
	ParseValue(value string) (any, error)
	Update(storage repositories.MetricsRepository, name string, value any) error
}

// Регистр
var processors = make(map[MetricType]MetricProcessor)

func RegisterProcessor(metricType MetricType, processor MetricProcessor) {
	processors[metricType] = processor
}

func GetProcessor(metricType MetricType) (MetricProcessor, error) {
	p, ok := processors[metricType]
	if !ok {
		return nil, customerrors.ErrUnsupportedType
	}
	return p, nil
}
