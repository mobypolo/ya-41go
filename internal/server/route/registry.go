package route

import "net/http"

type Route struct {
	Path    string
	Method  string
	Handler http.Handler
}

var routes []Route

func Register(path, method string, handler http.Handler) {
	routes = append(routes, Route{Path: path, Method: method, Handler: handler})
}

func All() []Route {
	return routes
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
