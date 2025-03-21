package storage

import (
	"github.com/mobypolo/ya-41go/internal/customerrors"
)

type MetricProcessor interface {
	ValidateName(name string) error
	ParseValue(value string) (any, error)
	Update(storage Storage, name string, value any) error
}

// Регистр
var processors = make(map[string]MetricProcessor)

func RegisterProcessor(metricType string, processor MetricProcessor) {
	processors[metricType] = processor
}

func GetProcessor(metricType string) (MetricProcessor, error) {
	p, ok := processors[metricType]
	if !ok {
		return nil, customerrors.ErrUnsupportedType
	}
	return p, nil
}
