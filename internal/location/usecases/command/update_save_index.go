package command

import (
	"context"
	"github.com/ahmedkamals/location_history/internal/location/domain/location"
)

type (
	UpdateSaveIndex struct {
		SaveIndex uint64
	}

	UpdateSaveIndexHandler struct {
		repo location.Repository
	}
)

func NewUpdateSaveIndexHandler(repo location.Repository) UpdateSaveIndexHandler {
	if repo == nil {
		panic("nil repo")
	}

	return UpdateSaveIndexHandler{
		repo: repo,
	}
}

func (a UpdateSaveIndexHandler) Execute(ctx context.Context, cmd UpdateSaveIndex) {
	defer func() {
		// Todo: log command execution.
	}()

	a.repo.UpdateSaveIndex(cmd.SaveIndex)
}
