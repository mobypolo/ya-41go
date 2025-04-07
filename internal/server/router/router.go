package router

import (
	"github.com/mobypolo/ya-41go/internal/server/route"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/server/handler"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	for _, r := range route.All() {
		mux.Handle(r.Path, r.Handler)
	}

	return mux
}
