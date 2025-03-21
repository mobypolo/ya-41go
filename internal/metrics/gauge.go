package metrics

import (
	"github.com/mobypolo/ya-41go/internal/customerrors"
	"strconv"

	"github.com/mobypolo/ya-41go/internal/storage"
)

type gaugeProcessor struct{}

func (g gaugeProcessor) ValidateName(name string) error {
	_, err := storage.ParseGaugeMetric(name)
	return err
}

func (g gaugeProcessor) ParseValue(value string) (any, error) {
	return strconv.ParseFloat(value, 64)
}

func (g gaugeProcessor) Update(s storage.Storage, name string, value any) error {
	v, ok := value.(float64)
	if !ok {
		return customerrors.ErrInvalidValue
	}
	return s.UpdateGauge(name, v)
}
