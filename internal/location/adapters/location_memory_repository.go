package adapters

import (
	"context"
	"github.com/ahmedkamals/location_history/internal/errors"
	"github.com/ahmedkamals/location_history/internal/location/domain/location"
	"github.com/ahmedkamals/location_history/internal/location/usecases/query"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type (
	LocationModel struct {
		orderUUID string
		latitude  float64
		longitude float64
		createdAt time.Time
	}

	LocationMemoryRepository struct {
		saveIndex uint64
		data      sync.Map
	}
)

func NewLocationMemoryRepository() *LocationMemoryRepository {
	return &LocationMemoryRepository{}
}

func (l *LocationMemoryRepository) AddLocation(ctx context.Context, order location.Order, lat, lng float64, createdAt time.Time) error {
	locations, ok := l.data.Load(order.OrderUUID)

	if !ok {
		locations = []LocationModel{}
	}

	locations = append([]LocationModel{{
		orderUUID: order.OrderUUID,
		latitude:  lat,
		longitude: lng,
		createdAt: createdAt,
	}}, locations.([]LocationModel)...)

	l.data.Store(order.OrderUUID, locations)

	return nil
}

func (l *LocationMemoryRepository) FindLocationsForOrder(ctx context.Context, criteria query.LocationForOrderCriteria) (query.History, error) {
	const op errors.Operation = "LocationMemoryRepository.FindLocationsForOrder"
	locations, ok := l.data.Load(criteria.OrderUUID)

	if !ok {
		emptyHistory := NewLocationsHistory(criteria.OrderUUID, []query.Location{})

		return emptyHistory, errors.E(op, errors.NotFound)
	}

	locationsData := locations.([]LocationModel)
	if criteria.Limit == nil {
		return l.locationsForOrder(criteria.OrderUUID, locationsData)
	}

	size := uint64(len(locationsData))

	if *criteria.Limit > size {
		*criteria.Limit = size
	}

	return l.locationsForOrder(criteria.OrderUUID, locationsData[:*criteria.Limit])
}

func (l *LocationMemoryRepository) locationsForOrder(orderUUID string, locations []LocationModel) (query.History, error) {
	matchedLocations := make([]query.Location, len(locations))

	for index, loc := range locations {
		matchedLocations[index] = query.Location{
			Lat: loc.latitude,
			Lon: loc.longitude,
		}
	}

	return NewLocationsHistory(orderUUID, matchedLocations), nil
}

func (l *LocationMemoryRepository) RemoveLocation(order location.Order) {
	l.data.Delete(order.OrderUUID)
}

func (l *LocationMemoryRepository) UpdateSaveIndex(saveIndex uint64) {
	var mux sync.Mutex

	mux.Lock()
	defer mux.Unlock()

	l.saveIndex = saveIndex
}

func (l *LocationMemoryRepository) RemoveAllExpiredLocations() {
	l.data.Range(func(key, value interface{}) bool {
		logrus.Info(key, value)
		return true
	})
}
