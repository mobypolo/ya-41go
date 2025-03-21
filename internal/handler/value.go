package handler

import (
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/customerrors"
	"github.com/mobypolo/ya-41go/internal/helpers"
	"github.com/mobypolo/ya-41go/internal/route"
	"github.com/mobypolo/ya-41go/internal/service"
	"github.com/mobypolo/ya-41go/internal/storage"
	"log"
	"net/http"
)

func init() {
	metricService := service.NewMetricService(storage.NewMemStorage())
	route.Register("/value/", MakeValueHandler(metricService))
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
		_, err = fmt.Fprint(w, val)
		if err != nil {
			log.Println(customerrors.ErrNotFound)
		}
	}
}
