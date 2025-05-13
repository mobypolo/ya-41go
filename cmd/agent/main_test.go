//go:build agent
// +build agent

package main

import (
	"fmt"
	"github.com/mobypolo/ya-41go/cmd"
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/stretchr/testify/assert"
)

func TestSendMetric_OK(t *testing.T) {
	var received string

	// Мокаем HTTP-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.URL.Path
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Подменяем адрес сервера
	originalAddress := cmd.ServerAddress
	cmd.ServerAddress = strings.TrimPrefix(ts.URL, "http://")
	defer func() { cmd.ServerAddress = originalAddress }()

	m := agent.Metric{
		Name:  "TestMetric",
		Type:  storage.GaugeType,
		Value: 42.42,
	}

	sendMetric(m)

	expectedPath := fmt.Sprintf("/update_plain/%s/%s/%v", m.Type, m.Name, m.Value)
	assert.Equal(t, expectedPath, received)
}

func TestSendMetric_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	// Подменяем адрес сервера
	originalAddress := cmd.ServerAddress
	cmd.ServerAddress = strings.TrimPrefix(ts.URL, "http://")
	defer func() { cmd.ServerAddress = originalAddress }()

	m := agent.Metric{
		Name:  "ErrorMetric",
		Type:  storage.CounterType,
		Value: 1,
	}

	sendMetric(m)
}
