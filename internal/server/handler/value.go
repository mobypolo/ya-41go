package handler

import (
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/server/helpers"
	"github.com/mobypolo/ya-41go/internal/server/route"
	"github.com/mobypolo/ya-41go/internal/server/service"
	"log"
	"net/http"
)

func init() {
	route.Register("/value/*", http.MethodGet, MakeValueHandler(service.GetMetricService()))
}

func MakeValueHandler(service *service.MetricService) http.Handler {
	var h http.Handler = ValueHandler(service)

	return h
}

func ValueHandler(service *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := helpers.SplitStrToChunks(r.URL.Path)

		metricType := parts[1]
		metricName := parts[2]

		val, err := service.Get(metricType, metricName)
		if err != nil {
			switch {
			case errors.Is(err, customerrors.ErrNotFound):
				http.Error(w, "metric not found", http.StatusNotFound)
			case errors.Is(err, customerrors.ErrUnsupportedType):
				http.Error(w, err.Error(), http.StatusNotImplemented)
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(w, "%s\n", val)
		if err != nil {
			log.Println(customerrors.ErrNotFound)
		}
	}
}
