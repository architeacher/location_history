package usecases

import (
	"context"
	"github.com/ahmedkamals/location_history/internal/location/usecases/command"
	"github.com/ahmedkamals/location_history/internal/location/usecases/query"
	"sync/atomic"
	"time"
)

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct {
		AddLocationHandler               command.AddLocationHandler
		RemoveLocationHandler            command.RemoveLocationHandler
		UpdateSaveIndexHandler           command.UpdateSaveIndexHandler
		RemoveAllExpiredLocationsHandler command.RemoveAllExpiredLocationsHandler
	}

	Queries struct {
		AllLocationsForOrderHandler       query.AllLocationForOrderHandler
		LocationsForOrderWithLimitHandler query.LocationForOrderWithLimitHandler
	}
)

func (app Application) EnableAutomaticCleanup(ctx context.Context, ttl time.Duration) {
	saveIndex := uint64(0)

	for range time.Tick(time.Second) {
		atomic.AddUint64(&saveIndex, 1)
		saveIndex %= uint64(ttl.Seconds())

		app.Commands.UpdateSaveIndexHandler.Execute(ctx, command.UpdateSaveIndex{SaveIndex: saveIndex})

		//app.Commands.RemoveAllExpiredLocationsHandler.Execute(
		//	context.Background(),
		//	command.RemoveAllExpiredLocations{
		//		TTL: ttl,
		//	})
	}
}
