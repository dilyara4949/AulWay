package domain

import "time"

type FavoriteRoute struct {
	UserID    string    `json:"user_id"`
	RouteID   string    `json:"route_id"`
	CreatedAt time.Time `json:"created_at"`
}
