package router

import (
	"net/http"
)

func MakeRouteHandler(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middlewares {
		h = mw(h)
	}
	return h
}
