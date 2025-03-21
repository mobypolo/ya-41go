package handler

import (
	"errors"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/customerrors"
	"github.com/mobypolo/ya-41go/internal/helpers"
	"github.com/mobypolo/ya-41go/internal/middleware"
	"github.com/mobypolo/ya-41go/internal/route"
	"github.com/mobypolo/ya-41go/internal/service"
	"github.com/mobypolo/ya-41go/internal/storage"
	"log"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/metrics"

func init() {
	metricService := service.NewMetricService(storage.NewMemStorage())
	route.Register("/update/", MakeUpdateHandler(metricService))
}

func MakeUpdateHandler(service *service.MetricService) http.Handler {
	var h http.Handler = UpdateHandler(service)

	h = middleware.RequirePathParts(4, h)
	h = middleware.AllowOnlyPost(h)

	return h
}

func UpdateHandler(service *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		parts := helpers.SplitStrToChunks(r.URL.Path)

		metricType, metricName, metricValue := parts[1], parts[2], parts[3]

		if err := service.Update(metricType, metricName, metricValue); err != nil {
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
}
