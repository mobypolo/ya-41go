package handler

import (
	"encoding/json"
	"fmt"
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/server/helpers"
	"github.com/mobypolo/ya-41go/internal/server/middleware"
	"github.com/mobypolo/ya-41go/internal/server/route"
	"github.com/mobypolo/ya-41go/internal/server/router"
	"github.com/mobypolo/ya-41go/internal/server/service"
	"github.com/mobypolo/ya-41go/internal/shared/dto"
	"log"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/server/metrics"

func init() {
	route.Register("/update/*", http.MethodPost, router.MakeRouteHandler(UpdateHandler(service.GetMetricService()), middleware.AllowOnlyPost, middleware.RequirePathParts(4)))
	route.Register("/update", http.MethodPost, router.MakeRouteHandler(UpdateJSONHandler(service.GetMetricService()), middleware.AllowOnlyPost, middleware.SetJSONContentType))
}

func UpdateHandler(service *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		parts := helpers.SplitStrToChunks(r.URL.Path)

		metricType, metricName, metricValue := parts[1], parts[2], parts[3]

		if err := service.Update(metricType, metricName, metricValue); err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		_, err := fmt.Fprintf(w, "Metric %s/%s updated with value %s\n", metricType, metricName, metricValue)
		if err != nil {
			log.Println(customerrors.ErrNotFound)
		}
	}
}

func UpdateJSONHandler(service *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m dto.Metrics
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if err := service.UpdateFromDTO(m); err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		actual, err := service.GetAsDTO(m.MType, m.ID)
		if err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(actual)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
	}
}
