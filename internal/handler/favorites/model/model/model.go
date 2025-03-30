package model

import "time"

type AddFavoriteRequest struct {
	RouteID string `json:"route_id" validate:"required"`
}

type RouteResponse struct {
	Id             string    `json:"id"`
	Departure      string    `json:"departure"`
	Destination    string    `json:"destination"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	AvailableSeats int       `json:"available_seats"`
	BusId          string    `json:"bus_id"`
	Price          int       `json:"price"`
	BusNumber      string    `json:"bus_number"`
	BusTotalSeats  int       `json:"bus_total_seats"`
}
