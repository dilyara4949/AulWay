package model

import (
	"errors"
)

type CreateRequest struct {
	Number     string `json:"number"`
	TotalSeats int    `json:"total_seats"`
}

func (createRequest CreateRequest) Validate() error {
	if createRequest.Number == "" || createRequest.TotalSeats <= 0 {
		return errors.New("required fields are cannot be empty or is zero")
	}
	return nil
}
