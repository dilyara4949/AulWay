package domain

import "time"

type Ticket struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	RouteID       string    `json:"route_id"`
	Price         int       `json:"price"`
	Status        string    `json:"status"`         // "approved", "cancelled", "awaiting"
	PaymentStatus string    `json:"payment_status"` // "pending", "paid", "failed"
	QRCode        string    `json:"qr_code"`
	CreatedAt     time.Time `json:"created_at"`
}
