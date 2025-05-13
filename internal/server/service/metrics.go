package service

import (
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/server/repositories"
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"github.com/mobypolo/ya-41go/internal/shared/dto"
	"strconv"
)

var (
	metricService *MetricService
)

func GetMetricService() *MetricService {
	return metricService
}

func SetMetricService(service *MetricService) {
	metricService = service
}

type MetricService struct {
	store repositories.MetricsRepository
}

func NewMetricService(store repositories.MetricsRepository) *MetricService {
	return &MetricService{store: store}
}

func (s *MetricService) Update(metricType storage.MetricType, name, value string) error {
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

func (s *MetricService) Get(metricType storage.MetricType, name string) (string, error) {
	switch metricType {
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

func (s *MetricService) UpdateFromDTO(m dto.Metrics) error {
	switch m.MType {
	case storage.GaugeType:
		if m.Value == nil {
			return customerrors.ErrInvalidValue
		}
		return s.store.UpdateGauge(m.ID, *m.Value)
	case storage.CounterType:
		if m.Delta == nil {
			return customerrors.ErrInvalidValue
		}
		return s.store.UpdateCounter(m.ID, *m.Delta)
	default:
		return customerrors.ErrUnsupportedType
	}
}

func (s *MetricService) GetAsDTO(mType storage.MetricType, id string) (dto.Metrics, error) {
	switch mType {
	case storage.GaugeType:
		val, err := s.store.GetGauge(id)
		if err != nil {
			return dto.Metrics{}, err
		}
		return dto.Metrics{ID: id, MType: mType, Value: &val}, nil
	case storage.CounterType:
		val, err := s.store.GetCounter(id)
		if err != nil {
			return dto.Metrics{}, err
		}
		return dto.Metrics{ID: id, MType: mType, Delta: &val}, nil
	default:
		return dto.Metrics{}, customerrors.ErrUnsupportedType
	}
}
