package handler

import (
	"fmt"
	"github.com/mobypolo/ya-41go/internal/server/route"
	"github.com/mobypolo/ya-41go/internal/server/service"
	"net/http"
	"sort"
)

func init() {
	route.Register("/", http.MethodPost, MakeIndexHandler(service.GetMetricService()))
}

func MakeIndexHandler(service *service.MetricService) http.Handler {
	var h http.Handler = IndexHandler(service)

	return h
}

func IndexHandler(service *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		_, err := fmt.Fprintf(w, "<html><body><h1>Metrics</h1><ul>")
		if err != nil {
			return
		}

		metrics := service.GetAvailableMetrics()

		keys := make([]string, 0, len(metrics))
		for k := range metrics {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			_, err := fmt.Fprintf(w, "<li>%s: %s</li>", k, metrics[k])
			if err != nil {
				return
			}
		}

		_, err = fmt.Fprintf(w, "</ul></body></html>")
		if err != nil {
			return
		}
	}
}
