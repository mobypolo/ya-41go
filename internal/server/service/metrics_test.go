//go:build server
// +build server

package service_test

import (
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/server/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

import _ "github.com/mobypolo/ya-41go/internal/server/metrics"

type mockRepo struct {
	gaugeValue    float64
	counterValue  int64
	getGaugeErr   error
	getCounterErr error
	updateCalled  bool
}

func (m *mockRepo) UpdateGauge(_ string, value float64) error {
	m.updateCalled = true
	m.gaugeValue = value
	return nil
}
func (m *mockRepo) UpdateCounter(_ string, delta int64) error {
	m.updateCalled = true
	m.counterValue += delta
	return nil
}
func (m *mockRepo) GetGauge(_ string) (float64, error) {
	return m.gaugeValue, m.getGaugeErr
}
func (m *mockRepo) GetCounter(_ string) (int64, error) {
	return m.counterValue, m.getCounterErr
}

func (m *mockRepo) GetAllCounters() map[string]string {
	return map[string]string{}
}

func TestMetricService_UpdateGauge_OK(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewMetricService(repo)

	err := svc.Update("gauge", "temperature", "36.6")
	assert.NoError(t, err)
	assert.True(t, repo.updateCalled)
	assert.Equal(t, 36.6, repo.gaugeValue)
}

func TestMetricService_Update_InvalidGaugeName(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewMetricService(repo)

	err := svc.Update("gauge", "invalid", "36.6")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid name")
}

func TestMetricService_Update_InvalidValue(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewMetricService(repo)

	err := svc.Update("gauge", "temperature", "abc")
	assert.ErrorIs(t, err, customerrors.ErrInvalidValue)
}

func TestMetricService_GetGauge_OK(t *testing.T) {
	repo := &mockRepo{gaugeValue: 42.5}
	svc := service.NewMetricService(repo)

	val, err := svc.Get("gauge", "temperature")
	assert.NoError(t, err)
	assert.Equal(t, "42.500", val)
}

func TestMetricService_GetGauge_NotFound(t *testing.T) {
	repo := &mockRepo{getGaugeErr: customerrors.ErrNotFound}
	svc := service.NewMetricService(repo)

	_, err := svc.Get("gauge", "temperature")
	assert.ErrorIs(t, err, customerrors.ErrNotFound)
}

func TestMetricService_GetCounter_OK(t *testing.T) {
	repo := &mockRepo{counterValue: 99}
	svc := service.NewMetricService(repo)

	val, err := svc.Get("counter", "requests")
	assert.NoError(t, err)
	assert.Equal(t, "99", val)
}

func TestMetricService_Get_InvalidType(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewMetricService(repo)

	_, err := svc.Get("invalid", "foo")
	assert.ErrorIs(t, err, customerrors.ErrUnsupportedType)
}
