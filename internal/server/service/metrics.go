package service

import (
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/server/repositories"
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"strconv"
	"sync"
)

var (
	metricService *MetricService
	once          sync.Once
)

func GetMetricService() *MetricService {
	once.Do(func() {
		setMetricService(NewMetricService(storage.NewMemStorage()))
	})
	return metricService
}

func setMetricService(service *MetricService) {
	metricService = service
}

type MetricService struct {
	store repositories.MetricsRepository
}

func NewMetricService(store repositories.MetricsRepository) *MetricService {
	return &MetricService{store: store}
}

func (s *MetricService) Update(metricType, name, value string) error {
	processor, err := storage.GetProcessor(metricType)
	if err != nil {
		return err
	}

	var allErrs []error

	if err := processor.ValidateName(name); err != nil {
		allErrs = append(allErrs, fmt.Errorf("invalid name: %w", err))
	}

	parsedValue, err := processor.ParseValue(value)
	if err != nil {
		allErrs = append(allErrs, customerrors.ErrInvalidValue)
	}

	if len(allErrs) > 0 {
		return errors.Join(allErrs...)
	}

	return processor.Update(s.store, name, parsedValue)
}

func (s *MetricService) Get(metricType, name string) (string, error) {
	switch storage.MetricType(metricType) {
	case storage.GaugeType:
		if err := s.validateGaugeName(name); err != nil {
			return "", err
		}
		val, err := s.store.GetGauge(name)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(val, 'f', -1, 64), nil

	case storage.CounterType:
		if err := s.validateCounterName(name); err != nil {
			return "", err
		}
		val, err := s.store.GetCounter(name)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", val), nil

	default:
		return "", customerrors.ErrUnsupportedType
	}
}

func (s *MetricService) GetAvailableMetrics() map[string]string {
	return s.store.GetAllCounters()
}

func (s *MetricService) validateGaugeName(name string) error {
	_, err := storage.ParseGaugeMetric(name)
	return err
}

func (s *MetricService) validateCounterName(name string) error {
	_, err := storage.ParseCounterMetric(name)
	return err
}
