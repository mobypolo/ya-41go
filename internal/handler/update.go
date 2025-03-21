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

	h = middleware.RequirePathParts(4, h)
	h = middleware.AllowOnlyPost(h)

	route.Register("/update/", h)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	parts := helpers.SplitStrToChunks(r.URL.Path)

	metricType, metricName, metricValue := parts[1], parts[2], parts[3]

	if err := memStore.UpdateMetric(metricType, metricName, metricValue); err != nil {
		switch {
		case errors.Is(err, customerrors.ErrUnsupportedType):
			http.Error(w, err.Error(), http.StatusNotImplemented)
		case errors.Is(err, customerrors.ErrUnknownGaugeName):
		case errors.Is(err, customerrors.ErrUnknownCounterName):
		case errors.Is(err, customerrors.ErrInvalidValue):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
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
