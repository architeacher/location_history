package location

import (
	"context"
	"time"
)

type (
	Repository interface {
		AddLocation(context.Context, Order, float64, float64, time.Time) error
		RemoveLocation(Order)
		UpdateSaveIndex(uint64)
		RemoveAllExpiredLocations()
	}
)
