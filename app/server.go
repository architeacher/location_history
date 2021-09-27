package app

import (
	"context"
	"errors"
	"github.com/ahmedkamals/location_history/config"
	"github.com/ahmedkamals/location_history/infrastructure"
	"github.com/ahmedkamals/location_history/internal/location/ports"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type (
	// HTTPServer is an interface for handling HTTP connections.
	HTTPServer struct {
		config config.Configuration
		routes []ports.Route
		server *http.Server
	}
)

// NewHTTPServer allocates and returns a new HTTPServer to handle HTTP connections.
func NewHTTPServer(
	config config.Configuration,
	routes []ports.Route,
) *HTTPServer {
	return &HTTPServer{
		config: config,
		routes: routes,
		server: &http.Server{
			Addr:         config.Port,
			WriteTimeout: config.WriteTimeout,
			ReadTimeout:  config.ReadTimeout * time.Second,
			IdleTimeout:  config.TTL,
		},
	}
}

// Start initiates routes configuration, and starts listening.
func (h *HTTPServer) Start() error {
	router, err := h.configureRoutesHandler(mux.NewRouter())

	if err != nil {
		return err
	}

	h.server.Handler = router

	return h.server.ListenAndServe()
}

func (h *HTTPServer) configureRoutesHandler(router *mux.Router) (*mux.Router, error) {
	if len(h.routes) == 0 {
		return nil, errors.New("missing routes")
	}

	for _, route := range h.routes {
		router.Handle(route.Path, route.Handler)
	}

	router.Use(infrastructure.LoggingMiddleware)

	return router, nil
}

func (h *HTTPServer) Shutdown(controlCtx context.Context) error {
	return h.server.Shutdown(controlCtx)
}

// BlockingClose closes all active connections gracefully.
func (h *HTTPServer) BlockingClose() {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	<-c
}
