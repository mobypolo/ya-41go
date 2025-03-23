//go:build server
// +build server

package handler_test

import (
	"github.com/mobypolo/ya-41go/internal/server/handler"
	"github.com/mobypolo/ya-41go/internal/server/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	GetGaugeFunc   func(name string) (float64, error)
	GetCounterFunc func(name string) (int64, error)
}

func (m *mockStore) GetGauge(name string) (float64, error) {
	return m.GetGaugeFunc(name)
}

func (m *mockStore) GetCounter(name string) (int64, error) {
	return m.GetCounterFunc(name)
}

func (m *mockStore) UpdateGauge(_ string, _ float64) error {
	return nil
}

func (m *mockStore) UpdateCounter(_ string, _ int64) error {
	return nil
}

func (m *mockStore) GetAllCounters() map[string]string {
	return map[string]string{}
}

func TestValueHandler_OK(t *testing.T) {
	s := &mockStore{
		GetGaugeFunc: func(name string) (float64, error) {
			assert.Equal(t, "temperature", name)
			return 123.45, nil
		},
	}

	svc := service.NewMetricService(s)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/temperature", nil)
	rr := httptest.NewRecorder()

	h := handler.ValueHandler(svc)
	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "123.450\n", rr.Body.String())
}

func TestValueHandler_UnsupportedType(t *testing.T) {
	s := &mockStore{
		GetGaugeFunc: func(name string) (float64, error) {
			assert.Equal(t, "temperature", name)
			return 42.42, nil
		},
	}

	svc := service.NewMetricService(s)

	req := httptest.NewRequest(http.MethodGet, "/value/invalid/metric", nil)
	rr := httptest.NewRecorder()

	h := handler.ValueHandler(svc)
	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotImplemented, rr.Code)
}

func TestValueHandler_BadRequest(t *testing.T) {
	s := &mockStore{
		GetGaugeFunc: func(name string) (float64, error) {
			assert.Equal(t, "temperature", name)
			return 42.42, nil
		},
	}

	svc := service.NewMetricService(s)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/broken", nil)
	rr := httptest.NewRecorder()

	h := handler.ValueHandler(svc)
	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
