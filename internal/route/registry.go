package route

import "net/http"

type Route struct {
	Path    string
	Handler http.Handler
}

var routes []Route

func Register(path string, handler http.Handler) {
	routes = append(routes, Route{Path: path, Handler: handler})
}

func All() []Route {
	return routes
}
