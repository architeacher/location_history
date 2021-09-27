package query

import (
	"context"
)

type (
	LocationForOrderCriteria struct {
		OrderUUID string
		Limit     *uint64
	}

	History interface {
		GetOrderUUID() string
		GetLocations() []Location
	}

	LocationForOrderReadModel interface {
		FindLocationsForOrder(context.Context, LocationForOrderCriteria) (History, error)
	}

	AllLocationForOrderHandler struct {
		readModel LocationForOrderReadModel
	}

	LocationForOrderWithLimitHandler struct {
		readModel LocationForOrderReadModel
	}
)

func NewLocationForOrderHandler(readModel LocationForOrderReadModel) AllLocationForOrderHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return AllLocationForOrderHandler{readModel: readModel}
}

func NewLocationForOrderWithLimitHandler(readModel LocationForOrderReadModel) LocationForOrderWithLimitHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return LocationForOrderWithLimitHandler{readModel: readModel}
}

func (l AllLocationForOrderHandler) Execute(ctx context.Context, orderUUID string) (history History, err error) {
	return l.readModel.FindLocationsForOrder(ctx, LocationForOrderCriteria{OrderUUID: orderUUID})
}

func (l LocationForOrderWithLimitHandler) Execute(ctx context.Context, orderUUID string, limit uint64) (history History, err error) {
	return l.readModel.FindLocationsForOrder(ctx, LocationForOrderCriteria{OrderUUID: orderUUID, Limit: &limit})
}
