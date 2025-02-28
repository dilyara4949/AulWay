package errs

import "errors"

var (
	ErrRecordNotFound       = errors.New("record not found")
	EmailAlreadyExists      = errors.New("email already exists")
	ErrInvalidEmailPassword = "invalid email or password"
	ErrTicketOutOfStock     = "tickets are out of stock"
)
