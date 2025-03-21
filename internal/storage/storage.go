package storage

import (
	"github.com/mobypolo/ya-41go/internal/customErrors"
	"sync"
)

type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

type Storage interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, delta int64) error
	GetGauge(string) (float64, error)
	GetCounter(string) (int64, error)
}

type MemStorage struct {
	mu       sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateMetric(metricType, name, value string) error {
	processor, err := GetProcessor(metricType)
	if err != nil {
		return err
	}
	if err := processor.ValidateName(name); err != nil {
		return err
	}
	parsedValue, err := processor.ParseValue(value)
	if err != nil {
		return customErrors.ErrInvalidValue
	}
	return processor.Update(s, name, parsedValue)
}

func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Gauges[name] = value
	return nil
}

func (s *MemStorage) UpdateCounter(name string, delta int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Counters[name] += delta
	return nil
}

func (s *MemStorage) GetGauge(name string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.Gauges[name]
	if !ok {
		return 0, customErrors.ErrNotFound
	}
	return val, nil
}

func (s *MemStorage) GetCounter(name string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.Counters[name]
	if !ok {
		return 0, customErrors.ErrNotFound
	}
	return val, nil
}
