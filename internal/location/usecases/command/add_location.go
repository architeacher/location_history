package command

import (
	"context"
	"github.com/ahmedkamals/location_history/internal/location/domain/location"
	"time"
)

type (
	AddLocation struct {
		Order     location.Order
		Latitude  float64
		Longitude float64
		CreatedAt time.Time
	}

	AddLocationHandler struct {
		repo location.Repository
	}
)

func NewAddLocationHandler(repo location.Repository) AddLocationHandler {
	if repo == nil {
		panic("nil repo")
	}

	return AddLocationHandler{
		repo: repo,
	}
}

func (a AddLocationHandler) Execute(ctx context.Context, commands ...AddLocation) (err error) {
	defer func() {
		// Todo: log command execution.
	}()

	for _, cmd := range commands {
		err = a.repo.AddLocation(ctx, cmd.Order, cmd.Latitude, cmd.Longitude, cmd.CreatedAt)

		if err != nil {
			return
		}
	}

	return
}
