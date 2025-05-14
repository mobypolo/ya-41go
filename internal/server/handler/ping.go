package handler

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mobypolo/ya-41go/internal/server/service"
	"net/http"
	"time"
)

func PingHandler(_ service.MetricService, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "no DB configured", http.StatusInternalServerError)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		err := db.Ping(ctx)
		if err != nil {
			http.Error(w, "database unavailable", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	}
}
