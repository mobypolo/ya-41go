package metrics

import (
	"github.com/mobypolo/ya-41go/internal/customerrors"
	"github.com/mobypolo/ya-41go/internal/repositories"
	"strconv"

	"github.com/mobypolo/ya-41go/internal/storage"
)

type CounterProcessor struct{}

func NewCounterProcessor() *CounterProcessor {
	return &CounterProcessor{}
}

func (c CounterProcessor) ValidateName(name string) error {
	_, err := storage.ParseCounterMetric(name)
	return err
}

func (c CounterProcessor) ParseValue(value string) (any, error) {
	return strconv.ParseInt(value, 10, 64)
}

func (c CounterProcessor) Update(s repositories.MetricsRepository, name string, value any) error {
	v, ok := value.(int64)
	if !ok {
		return customerrors.ErrInvalidValue
	}
	return s.UpdateCounter(name, v)
}
