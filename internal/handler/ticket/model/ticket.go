package model

type BuyTicketRequest struct {
	Quantity  int    `json:"quantity"`
	UserEmail string `json:"user_email"`
}
