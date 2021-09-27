package service

import (
	"context"
	"github.com/ahmedkamals/location_history/config"
	"github.com/ahmedkamals/location_history/internal/location/adapters"
	"github.com/ahmedkamals/location_history/internal/location/usecases"
	"github.com/ahmedkamals/location_history/internal/location/usecases/command"
	"github.com/ahmedkamals/location_history/internal/location/usecases/query"
)

func NewApplication(ctx context.Context, config config.Configuration) usecases.Application {
	locationsRepository := adapters.NewLocationMemoryRepository()

	app := usecases.Application{
		Commands: usecases.Commands{
			AddLocationHandler:               command.NewAddLocationHandler(locationsRepository),
			RemoveLocationHandler:            command.NewRemoveLocationHandler(locationsRepository),
			UpdateSaveIndexHandler:           command.NewUpdateSaveIndexHandler(locationsRepository),
			RemoveAllExpiredLocationsHandler: command.NewRemoveAllExpiredLocationsHandler(locationsRepository),
		},
		Queries: usecases.Queries{
			AllLocationsForOrderHandler:       query.NewLocationForOrderHandler(locationsRepository),
			LocationsForOrderWithLimitHandler: query.NewLocationForOrderWithLimitHandler(locationsRepository),
		},
	}

	go app.EnableAutomaticCleanup(ctx, config.TTL)

	return app
}
