package handler

import (
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/customerrors"
	"github.com/mobypolo/ya-41go/internal/helpers"
	"github.com/mobypolo/ya-41go/internal/middleware"
	"github.com/mobypolo/ya-41go/internal/route"
	"github.com/mobypolo/ya-41go/internal/storage"
	"log"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/metrics"

var (
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
		case errors.Is(err, customerrors.ErrUnsupportedType):
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, err.Error(), http.StatusNotImplemented)
		case errors.Is(err, customerrors.ErrUnknownGaugeName):
		case errors.Is(err, customerrors.ErrUnknownCounterName):
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Metric %s/%s updated with value %s\n", metricType, metricName, metricValue)
	if err != nil {
		log.Println(customerrors.ErrNotFound)
	}
}
