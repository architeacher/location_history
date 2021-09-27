package ports

import (
	"github.com/ahmedkamals/location_history/internal/location/usecases"
	"net/http"
)

type (
	Route struct {
		Path    string
		Methods []string
		Scheme  string
		Handler http.Handler
	}
)

func NewRoute(path string, handler http.Handler, scheme string, methods ...string) Route {
	return Route{
		Path:    path,
		Methods: methods,
		Scheme:  scheme,
		Handler: handler,
	}
}

func GetRoutes(app usecases.Application, errQueue chan error) []Route {
	return []Route{
		NewRoute(
			"/location/{orderUUID}/now",
			NewHttpController(app, errQueue),
			"http",
			methodPost,
		),
		NewRoute(
			"/location/{orderUUID}",
			NewHttpController(app, errQueue),
			"http",
			methodGet,
		),
		NewRoute(
			"/location/{orderUUID}",
			NewHttpController(app, errQueue),
			"http",
			methodDelete,
		),
	}
}
