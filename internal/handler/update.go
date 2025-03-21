package handler

import (
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/customErrors"
	"github.com/mobypolo/ya-41go/internal/helpers"
	"github.com/mobypolo/ya-41go/internal/middleware"
	"github.com/mobypolo/ya-41go/internal/route"
	"github.com/mobypolo/ya-41go/internal/storage"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/metrics"

var (
	allowedGaugeMetrics   = map[string]struct{}{"temperature": {}, "load": {}}
	allowedCounterMetrics = map[string]struct{}{"requests": {}, "errors": {}}

	memStore = storage.NewMemStorage()
)

func init() {
	var h http.Handler = http.HandlerFunc(updateHandler)

	h = middleware.RequirePathParts(3, h) // /{type}/{name}/{value} — 3 части
	h = middleware.AllowOnlyPost(h)       // только POST

	route.Register("/update/", h)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	parts := helpers.SplitStrToChunks(r.URL.Path)

	metricType, metricName, metricValue := parts[1], parts[2], parts[3]

	if err := memStore.UpdateMetric(metricType, metricName, metricValue); err != nil {
		switch {
		case errors.Is(err, customErrors.ErrUnsupportedType):
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, err.Error(), http.StatusNotImplemented)
		case errors.Is(err, customErrors.ErrUnknownGaugeName):
		case errors.Is(err, customErrors.ErrUnknownCounterName):
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Metric %s/%s updated with value %s\n", metricType, metricName, metricValue)
}
