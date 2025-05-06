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
	"io"
	"log"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/server/metrics"

func init() {
	route.DeferRegister(func() {
		s := service.GetMetricService()
		if s == nil {
			panic("metricService not set before route registration")
		}
		route.Register("/update/*", http.MethodPost, router.MakeRouteHandler(UpdateHandler(service.GetMetricService()), middleware.AllowOnlyPost, middleware.RequirePathParts(4)))
		route.Register("/update/", http.MethodPost, router.MakeRouteHandler(UpdateJSONHandler(service.GetMetricService()), middleware.AllowOnlyPost, middleware.SetJSONContentType))
		route.Register("/updates/", http.MethodPost, router.MakeRouteHandler(UpdateJSONHandlerBatch(service.GetMetricService()), middleware.AllowOnlyPost, middleware.SetJSONContentType))
	})
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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		var single dto.Metrics
		if err := json.Unmarshal(body, &single); err == nil {
			if err := service.UpdateFromDTO(single); err != nil {
				customerrors.ErrorHandler(err, w)
				return
			}

			actual, err := service.GetAsDTO(single.MType, single.ID)
			if err != nil {
				customerrors.ErrorHandler(err, w)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(actual)
			if err != nil {
				return
			}
			return
		}

		http.Error(w, "invalid JSON format", http.StatusBadRequest)
	}
}
func UpdateJSONHandlerBatch(service *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		var batch []dto.Metrics
		if err := json.Unmarshal(body, &batch); err == nil {
			for _, metric := range batch {
				if err := service.UpdateFromDTO(metric); err != nil {
					customerrors.ErrorHandler(err, w)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		http.Error(w, "invalid JSON format", http.StatusBadRequest)
	}
}
