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
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"github.com/mobypolo/ya-41go/internal/shared/dto"
	"log"
	"net/http"
)

func init() {
	route.DeferRegister(func() {
		s := service.GetMetricService()
		if s == nil {
			panic("metricService not set before route registration")
		}
		route.Register("/value/*", http.MethodGet, router.MakeRouteHandler(ValueHandler(service.GetMetricService())))
		route.Register("/value/", http.MethodPost, router.MakeRouteHandler(ValueJSONHandler(service.GetMetricService()), middleware.SetJSONContentType))
	})
}

func ValueHandler(service *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := helpers.SplitStrToChunks(r.URL.Path)

		metricType := storage.MetricType(parts[1])
		metricName := parts[2]

		val, err := service.Get(metricType, metricName)
		if err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(w, "%s\n", val)
		if err != nil {
			log.Println(customerrors.ErrNotFound)
		}
	}
}

func ValueJSONHandler(service *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.Metrics
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		resp, err := service.GetAsDTO(req.MType, req.ID)
		if err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
	}
}
