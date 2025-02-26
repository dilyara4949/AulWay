package errs

import "errors"

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrInvalidEmailPassword = "invalid email or password"
	ErrTicketOutOfStock     = "tickets are out of stock"
)
