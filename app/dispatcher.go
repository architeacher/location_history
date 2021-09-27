package app

import (
	"context"
	"github.com/ahmedkamals/location_history/config"
	"github.com/ahmedkamals/location_history/infrastructure"
	"github.com/ahmedkamals/location_history/internal/location/ports"
	"github.com/ahmedkamals/location_history/internal/location/service"
)

type (
	// Dispatcher of the application.
	Dispatcher struct {
		config   config.Configuration
		errQueue chan error
		logQueue chan string
	}
)

// NewDispatcher creates a new Dispatcher instance.
func NewDispatcher(config config.Configuration) *Dispatcher {
	return &Dispatcher{
		config:   config,
		errQueue: make(chan error, config.BufferSize),
	}
}

// Dispatch operation.
func (d Dispatcher) Dispatch(controlCtx context.Context) {
	app := service.NewApplication(controlCtx, d.config)

	infrastructure.LogQueue <- "Starting server..."
	server := NewHTTPServer(d.config, ports.GetRoutes(app, d.errQueue))

	go infrastructure.MonitorLogMessages(infrastructure.LogQueue)
	go infrastructure.MonitorRequestLogMessages(infrastructure.ReqLogQueue)
	go infrastructure.MonitorErrors(d.errQueue)

	go func() {
		d.errQueue <- server.Start()
	}()

	server.BlockingClose()
}
