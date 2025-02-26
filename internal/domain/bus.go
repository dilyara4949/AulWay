package domain

type Bus struct {
	Id         string `json:"id"`
	Number     string `json:"number"`
	TotalSeats int    `json:"total_seats"`
}
