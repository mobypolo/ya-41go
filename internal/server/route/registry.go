package route

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mobypolo/ya-41go/cmd"
	"github.com/mobypolo/ya-41go/internal/server/handler"
	"github.com/mobypolo/ya-41go/internal/server/middleware"
	"github.com/mobypolo/ya-41go/internal/server/router"
	"github.com/mobypolo/ya-41go/internal/server/service"
	"net/http"
)

type Route struct {
	Path    string
	Method  string
	Handler http.Handler
}

var routes []Route

func Register(path, method string, handler http.Handler) {
	routes = append(routes, Route{Path: path, Method: method, Handler: handler})
}

func RegisterAllRoutes(db *pgxpool.Pool, cfg cmd.Config) {
	s := service.GetMetricService()
	if s == nil {
		panic("metricService not set before route registration")
	}
	Register("/", http.MethodPost, router.MakeRouteHandler(handler.IndexHandler(s), middleware.HashSHA256(cfg.Key)))
	Register("/ping", http.MethodGet, router.MakeRouteHandler(handler.PingHandler(s, db)))
	Register("/update/*", http.MethodPost, router.MakeRouteHandler(handler.UpdateHandler(s), middleware.HashSHA256(cfg.Key), middleware.AllowOnlyPost, middleware.RequirePathParts(4)))
	Register("/update/", http.MethodPost, router.MakeRouteHandler(handler.UpdateJSONHandler(s), middleware.HashSHA256(cfg.Key), middleware.AllowOnlyPost, middleware.SetJSONContentType))
	Register("/updates/", http.MethodPost, router.MakeRouteHandler(handler.UpdateJSONHandlerBatch(s), middleware.HashSHA256(cfg.Key), middleware.AllowOnlyPost, middleware.SetJSONContentType))
	Register("/value/*", http.MethodGet, router.MakeRouteHandler(handler.ValueHandler(s)))
	Register("/value/", http.MethodPost, router.MakeRouteHandler(handler.ValueJSONHandler(s), middleware.SetJSONContentType))
}

func MountInto(mux interface{}) {
	for _, r := range routes {
		switch m := mux.(type) {
		case interface {
			Method(method string, pattern string, h http.HandlerFunc)
		}:
			m.Method(r.Method, r.Path, r.Handler.ServeHTTP)
		case interface {
			Handle(pattern string, h http.Handler)
		}:
			m.Handle(r.Path, r.Handler)
		default:
			panic("unknown router")
		}
	}
}
