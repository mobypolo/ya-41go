package handler

import (
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/customErrors"
	"github.com/mobypolo/ya-41go/internal/helpers"
	"github.com/mobypolo/ya-41go/internal/route"
	"github.com/mobypolo/ya-41go/internal/service"
	"net/http"
)

var (
	metricService = service.NewMetricService(memStore)
)

func init() {
	route.Register("/value/", http.HandlerFunc(valueHandler))
}

func valueHandler(w http.ResponseWriter, r *http.Request) {
	parts := helpers.SplitStrToChunks(r.URL.Path)

	metricType := parts[1]
	metricName := parts[2]

	val, err := metricService.Get(metricType, metricName)
	if err != nil {
		switch {
		case errors.Is(err, customErrors.ErrNotFound):
			http.Error(w, "metric not found", http.StatusNotFound)
		case errors.Is(err, customErrors.ErrUnsupportedType):
			http.Error(w, err.Error(), http.StatusNotImplemented)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, val)
}
