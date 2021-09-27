package command

import (
	"context"
	"github.com/ahmedkamals/location_history/internal/location/domain/location"
)

type (
	RemoveLocation struct {
		Order location.Order
	}

	RemoveLocationHandler struct {
		repo location.Repository
	}
)

func NewRemoveLocationHandler(repo location.Repository) RemoveLocationHandler {
	if repo == nil {
		panic("nil repo")
	}

	return RemoveLocationHandler{
		repo: repo,
	}
}

func (a RemoveLocationHandler) Execute(ctx context.Context, commands ...RemoveLocation) {
	defer func() {
		// Todo: log command execution.
	}()

	for _, cmd := range commands {
		a.repo.RemoveLocation(cmd.Order)
	}
}
