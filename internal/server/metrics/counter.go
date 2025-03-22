package metrics

import (
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/server/repositories"
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"strconv"
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
