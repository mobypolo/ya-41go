package metrics

import (
	"github.com/mobypolo/ya-41go/internal/customerrors"
	"github.com/mobypolo/ya-41go/internal/repositories"
	"strconv"

	"github.com/mobypolo/ya-41go/internal/storage"
)

type GaugeProcessor struct{}

func NewGaugeProcessor() *GaugeProcessor {
	return &GaugeProcessor{}
}

func (g GaugeProcessor) ValidateName(name string) error {
	_, err := storage.ParseGaugeMetric(name)
	return err
}

func (g GaugeProcessor) ParseValue(value string) (any, error) {
	return strconv.ParseFloat(value, 64)
}

func (g GaugeProcessor) Update(s repositories.MetricsRepository, name string, value any) error {
	v, ok := value.(float64)
	if !ok {
		return customerrors.ErrInvalidValue
	}
	return s.UpdateGauge(name, v)
}
