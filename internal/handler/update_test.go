package handler_test

import (
	"github.com/mobypolo/ya-41go/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mobypolo/ya-41go/internal/handler"
	"github.com/mobypolo/ya-41go/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandler_OK(t *testing.T) {
	mem := storage.NewMemStorage()
	metricService := service.NewMetricService(mem)
	h := handler.MakeUpdateHandler(metricService)

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/temperature/36.6", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	val, err := mem.GetGauge("temperature")
	assert.NoError(t, err)
	assert.Equal(t, 36.6, val)
}

func TestUpdateHandler_BadMethod(t *testing.T) {
	mem := storage.NewMemStorage()
	metricService := service.NewMetricService(mem)
	h := handler.MakeUpdateHandler(metricService)

	req := httptest.NewRequest(http.MethodGet, "/update/gauge/temperature/36.6", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
