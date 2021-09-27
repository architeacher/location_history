package ports

import (
	"encoding/json"
	"github.com/ahmedkamals/location_history/infrastructure"
	"github.com/ahmedkamals/location_history/internal/errors"
	"github.com/ahmedkamals/location_history/internal/location/domain/location"
	"github.com/ahmedkamals/location_history/internal/location/usecases"
	"github.com/ahmedkamals/location_history/internal/location/usecases/command"
	"github.com/ahmedkamals/location_history/internal/location/usecases/query"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	LocationsResponse struct {
		OrderUUID string           `json:"order_id"`
		History   []query.Location `json:"history"`
	}

	httpController struct {
		app      usecases.Application
		errQueue chan error
	}
)

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodDelete = "DELETE"
)

func (l *httpController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bodyLength := 0

	defer func() {
		infrastructure.LogRequest(infrastructure.ReqLogQueue, &infrastructure.RequestLogMessage{
			Format:     infrastructure.CommonLogFormat,
			Request:    r,
			StatusCode: http.StatusOK,
			BodyLength: bodyLength,
		})
	}()

	acceptedResponseType := strings.ToLower(r.Header.Get("Accept"))

	if acceptedResponseType == "" || !strings.HasPrefix(acceptedResponseType, "application") {
		acceptedResponseType = "application/json; charset=UTF-8"
	}

	w.Header().Set("Content-Type", acceptedResponseType)

	var err error
	var response LocationsResponse

	switch strings.ToUpper(r.Method) {
	case methodPost:
		if err := l.appendLocation(r); err != nil {
			l.errQueue <- err
		}
		return
	case methodGet:
		if response, err = l.getLocations(r); err != nil {
			l.errQueue <- err
		}
	case methodDelete:
		l.deleteLocation(r)
		return
	default:
		return
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		l.errQueue <- err
	}

	bodyLength, err = w.Write(responseBody)
	if err != nil {
		l.errQueue <- err
	}

	return
}

func (l *httpController) appendLocation(r *http.Request) error {
	const op errors.Operation = "httpController.appendLocation"

	var newLocation location.EntryPoint

	err := json.NewDecoder(r.Body).Decode(&newLocation)
	if err != nil {
		return errors.E(op, errors.Failure, err)
	}

	orderUUID := mux.Vars(r)["orderUUID"]

	err = l.app.Commands.AddLocationHandler.Execute(r.Context(), command.AddLocation{
		Order:     location.NewOrder(orderUUID),
		Latitude:  newLocation.Latitude,
		Longitude: newLocation.Longitude,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return errors.E(op, errors.Failure, err)
	}

	return nil
}

func (l *httpController) getLocations(r *http.Request) (response LocationsResponse, err error) {
	const op errors.Operation = "httpController.getLocations"

	orderUUID := mux.Vars(r)["orderUUID"]
	maxEntries, hasLimit := r.URL.Query()["max"]
	if !hasLimit {
		history, err := l.app.Queries.AllLocationsForOrderHandler.Execute(r.Context(), orderUUID)

		return LocationsResponse{
			OrderUUID: history.GetOrderUUID(),
			History:   history.GetLocations(),
		}, err
	}

	emptyResponse := LocationsResponse{
		OrderUUID: orderUUID,
		History:   []query.Location{},
	}

	limit, err := strconv.Atoi(maxEntries[0])
	if err != nil {
		return emptyResponse, errors.E(op, errors.Failure, err)
	}

	history, err := l.app.Queries.LocationsForOrderWithLimitHandler.Execute(r.Context(), orderUUID, uint64(limit))
	if err != nil {
		return emptyResponse, errors.E(op, errors.Failure, err)
	}

	return LocationsResponse{
		OrderUUID: history.GetOrderUUID(),
		History:   history.GetLocations(),
	}, nil
}

func (l *httpController) deleteLocation(r *http.Request) {
	const op errors.Operation = "httpController.deleteLocation"

	orderUUID := mux.Vars(r)["orderUUID"]

	l.app.Commands.RemoveLocationHandler.Execute(r.Context(), command.RemoveLocation{
		Order: location.NewOrder(orderUUID),
	})
}

// NewHttpController allocates and returns a new instance of httpController.
func NewHttpController(app usecases.Application, errQueue chan error) http.Handler {
	return &httpController{
		app:      app,
		errQueue: errQueue,
	}
}
