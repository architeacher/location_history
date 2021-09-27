package adapters

import "github.com/ahmedkamals/location_history/internal/location/usecases/query"

type (
	LocationsHistory struct {
		OrderUUID string
		Locations []query.Location
	}
)

func NewLocationsHistory(OrderUUID string, locations []query.Location) query.History {
	return LocationsHistory{
		OrderUUID: OrderUUID,
		Locations: locations,
	}
}

func (l LocationsHistory) GetOrderUUID() string {
	return l.OrderUUID
}

func (l LocationsHistory) GetLocations() []query.Location {
	return l.Locations
}
