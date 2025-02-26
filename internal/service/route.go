package service

import (
	"aulway/internal/domain"
	"aulway/internal/handler/route/model"
	"aulway/internal/repository/route"
	"context"
	"fmt"
	"github.com/google/uuid"
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
		Id:             routeId.String(),
		Departure:      request.Departure,
		Destination:    request.Destination,
		StartDate:      request.StartDate,
		EndDate:        request.EndDate,
		BusId:          request.BusId,
		Price:          request.Price,
		AvailableSeats: availableSeats,
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
		updates["departure"] = req.Departure
	}
	if req.Destination != "" {
		updates["destination"] = req.Destination
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

func (service *Route) GetRoutesList(ctx context.Context, departure, destination string, date time.Time, passengers, page, pageSize int) ([]domain.Route, int, error) {
	return service.repo.GetRoutesList(ctx, departure, destination, date, passengers, page, pageSize)
}
