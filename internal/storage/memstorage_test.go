package storage_test

import (
	"github.com/mobypolo/ya-41go/internal/customerrors"
	"testing"

	"github.com/mobypolo/ya-41go/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	s := storage.NewMemStorage()

	err := s.UpdateGauge("temperature", 36.6)
	assert.NoError(t, err)

	val, err := s.GetGauge("temperature")
	assert.NoError(t, err)
	assert.Equal(t, 36.6, val)
}

func TestMemStorage_GetGauge(t *testing.T) {
	s := storage.NewMemStorage()

	// сначала сохраняем метрику
	err := s.UpdateGauge("temperature", 36.6)
	assert.NoError(t, err)

	// потом читаем
	val, err := s.GetGauge("temperature")
	assert.NoError(t, err)
	assert.Equal(t, 36.6, val)
}

func TestMemStorage_GetGauge_NotFound(t *testing.T) {
	s := storage.NewMemStorage()

	_, err := s.GetGauge("nonexistent")
	assert.ErrorIs(t, err, customerrors.ErrNotFound)
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	s := storage.NewMemStorage()

	err := s.UpdateCounter("requests", 5)
	assert.NoError(t, err)

	val, err := s.GetCounter("requests")
	assert.NoError(t, err)
	assert.Equal(t, int64(5), val)
}

func TestMemStorage_GetCounter(t *testing.T) {
	s := storage.NewMemStorage()

	err := s.UpdateCounter("requests", 10)
	assert.NoError(t, err)

	val, err := s.GetCounter("requests")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), val)
}

func TestMemStorage_GetCounter_NotFound(t *testing.T) {
	s := storage.NewMemStorage()

	_, err := s.GetCounter("nonexistent")
	assert.ErrorIs(t, err, customerrors.ErrNotFound)
}
