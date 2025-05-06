package handler

import (
	"context"
	"github.com/mobypolo/ya-41go/internal/server/db"
	"github.com/mobypolo/ya-41go/internal/server/route"
	"github.com/mobypolo/ya-41go/internal/server/router"
	"github.com/mobypolo/ya-41go/internal/server/service"
	"net/http"
	"time"
)

func init() {
	route.DeferRegister(func() {
		s := service.GetMetricService()
		if s == nil {
			panic("metricService not set before route registration")
		}
		route.Register("/ping", http.MethodGet, router.MakeRouteHandler(PingHandler(service.GetMetricService())))
	})
}

func PingHandler(_ *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db.Pool == nil {
			http.Error(w, "no DB configured", http.StatusInternalServerError)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		err := db.Pool.Ping(ctx)
		if err != nil {
			http.Error(w, "database unavailable", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	}
}
