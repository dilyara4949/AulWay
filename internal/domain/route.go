package domain

import "time"

type Route struct {
	Id             string    `json:"id"`
	Departure      string    `json:"departure"`
	Destination    string    `json:"destination"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	AvailableSeats int       `json:"available_seats"`
	BusId          string    `json:"bus_id"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsFavorite     bool      `json:"is_favorite" gorm:"-"`
}
