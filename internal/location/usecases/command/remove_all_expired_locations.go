package command

import (
	"context"
	"github.com/ahmedkamals/location_history/internal/location/domain/location"
	"time"
)

type (
	RemoveAllExpiredLocations struct {
		TTL time.Duration
	}

	RemoveAllExpiredLocationsHandler struct {
		repo location.Repository
	}
)

func NewRemoveAllExpiredLocationsHandler(repo location.Repository) RemoveAllExpiredLocationsHandler {
	if repo == nil {
		panic("nil repo")
	}

	return RemoveAllExpiredLocationsHandler{
		repo: repo,
	}
}

func (a RemoveAllExpiredLocationsHandler) Execute(ctx context.Context, command RemoveAllExpiredLocations) {
	defer func() {
		// Todo: log command execution.
	}()

	a.repo.RemoveAllExpiredLocations()
}
