package model

import (
	"aulway/internal/domain"
	"github.com/go-playground/validator/v10"
	"time"
)

type CreateRouteRequest struct {
	Departure           string    `json:"departure" validate:"required"`
	Destination         string    `json:"destination" validate:"required"`
	DepartureLocation   string    `json:"departure_location"`
	DestinationLocation string    `json:"destination_location"`
	StartDate           time.Time `json:"start_date" validate:"required" example:"2025-12-12T13:00:00+05:00"`
	EndDate             time.Time `json:"end_date" validate:"required,gtfield=StartDate" example:"2025-12-12T13:00:00+05:00"`
	BusId               string    `json:"bus_id" validate:"required"`
	Price               int       `json:"price" validate:"required,gte=0"`
}

type UpdateRouteRequest struct {
	Departure           string    `json:"departure" validate:"required"`
	Destination         string    `json:"destination" validate:"required"`
	DepartureLocation   string    `json:"departure_location"`
	DestinationLocation string    `json:"destination_location"`
	StartDate           time.Time `json:"start_date" validate:"required" example:"2025-12-12T13:00:00+05:00"`
	EndDate             time.Time `json:"end_date" validate:"required,gtfield=StartDate" example:"2025-12-12T13:00:00+05:00"`
	BusId               string    `json:"bus_id" validate:"required"`
	Price               int       `json:"price" validate:"required,gte=0"`
}

func (r *CreateRouteRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type RouteResponse struct {
	Id                  string    `json:"id"`
	Departure           string    `json:"departure"`
	Destination         string    `json:"destination"`
	DepartureLocation   string    `json:"departure_location"`
	DestinationLocation string    `json:"destination_location"`
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	AvailableSeats      int       `json:"available_seats"`
	BusId               string    `json:"bus_id"`
	Price               int       `json:"price"`
	BusNumber           string    `json:"bus_number"`
	BusTotalSeats       int       `json:"bus_total_seats"`
}

func MapRouteResponse(route domain.Route, bus domain.Bus) *RouteResponse {
	return &RouteResponse{
		Id:                  route.Id,
		Departure:           route.Departure,
		Destination:         route.Destination,
		DepartureLocation:   route.DepartureLocation,
		DestinationLocation: route.DestinationLocation,
		StartDate:           route.StartDate,
		EndDate:             route.EndDate,
		AvailableSeats:      route.AvailableSeats,
		BusId:               route.BusId,
		Price:               route.Price,
		BusNumber:           bus.Number,
		BusTotalSeats:       bus.TotalSeats,
	}
}
