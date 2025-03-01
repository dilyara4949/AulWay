package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type CreateRouteRequest struct {
	Departure   string    `json:"departure" validate:"required"`
	Destination string    `json:"destination" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required" example:"2025-12-12T13:00:00+05:00"`
	EndDate     time.Time `json:"end_date" validate:"required,gtfield=StartDate" example:"2025-12-12T13:00:00+05:00"`
	BusId       string    `json:"bus_id" validate:"required"`
	Price       int       `json:"price" validate:"required,gte=0"`
}

type UpdateRouteRequest struct {
	Departure   string    `json:"departure" validate:"required"`
	Destination string    `json:"destination" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required" example:"2025-12-12T13:00:00+05:00"`
	EndDate     time.Time `json:"end_date" validate:"required,gtfield=StartDate" example:"2025-12-12T13:00:00+05:00"`
	BusId       string    `json:"bus_id" validate:"required"`
	Price       int       `json:"price" validate:"required,gte=0"`
}

func (r *CreateRouteRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
