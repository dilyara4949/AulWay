package service

import (
	"aulway/internal/domain"
	"aulway/internal/handler/route/model"
	"aulway/internal/repository/route"
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"time"
)

type Route struct {
	repo route.Repository
}

func NewRouteService(routeRepo route.Repository) *Route {
	return &Route{
		repo: routeRepo,
	}
}

func (service *Route) CreateRoute(ctx context.Context, request model.CreateRouteRequest, availableSeats int) (*domain.Route, error) {
	routeId, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate uuid error: %w", err)
	}

	response := &domain.Route{
		Id:                  routeId.String(),
		Departure:           CapitalizeFirst(request.Departure),
		Destination:         CapitalizeFirst(request.Destination),
		DestinationLocation: request.DestinationLocation,
		DepartureLocation:   request.DepartureLocation,
		StartDate:           request.StartDate,
		EndDate:             request.EndDate,
		BusId:               request.BusId,
		Price:               request.Price,
		AvailableSeats:      availableSeats,
	}

	err = service.repo.Create(ctx, response)
	return response, err
}

func (service *Route) GetRoute(ctx context.Context, id string) (*domain.Route, error) {
	return service.repo.Get(ctx, id)
}

func (service *Route) Delete(ctx context.Context, id string) error {
	return service.repo.Delete(ctx, id)
}

func (service *Route) Update(ctx context.Context, req model.UpdateRouteRequest, id string) error {

	updates := make(map[string]interface{})

	if req.Departure != "" {
		updates["departure"] = CapitalizeFirst(req.Departure)
	}
	if req.Destination != "" {
		updates["destination"] = CapitalizeFirst(req.Destination)
	}
	if req.DepartureLocation != "" {
		updates["departure_location"] = CapitalizeFirst(req.DepartureLocation)
	}
	if req.Destination != "" {
		updates["destination_location"] = CapitalizeFirst(req.DestinationLocation)
	}
	if !req.StartDate.IsZero() {
		updates["start_date"] = req.StartDate
	}
	if !req.EndDate.IsZero() {
		updates["end_date"] = req.EndDate
	}
	if req.BusId != "" {
		updates["bus_id"] = req.BusId
	}
	if req.Price >= 0 {
		updates["price"] = req.Price
	}

	if len(updates) == 0 {
		return nil
	}

	return service.repo.Update(ctx, updates, id)
}

func (service *Route) GetRoutesListt(ctx context.Context, userId, departure, destination string, date time.Time, passengers, page, pageSize int) ([]domain.Route, int, error) {
	departure = CapitalizeFirst(departure)
	destination = CapitalizeFirst(destination)
	return service.repo.GetRoutesList(ctx, userId, departure, destination, date, passengers, page, pageSize)
}

func (service *Route) GetAllRoutesList(ctx context.Context, page, pageSize int) ([]domain.Route, error) {
	return service.repo.GetAllRoutesList(ctx, page, pageSize)
}

func CapitalizeFirst(str string) string {
	c := cases.Title(language.Und)
	return c.String(str)
}
