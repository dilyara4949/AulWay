package errs

import (
	"errors"
	"fmt"
)

type Err struct {
	ErrDesc string `json:"errDesc"`
	Err     string `json:"error"`
}

func (e Err) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.ErrDesc)
}

var ErrNoSeatsAvailable = errors.New("no seats available")
var ErrEmptyRequestFields = errors.New("request fields cannot be empty")
var ErrRequestBinding = errors.New("request binding error")
var ErrIncorrectPhoneFormat = errors.New("incorrect phone format error")
var ErrIncorrectEmailFormat = errors.New("incorrect email format error")
var ErrIncorrectPasswordFormat = errors.New("incorrect password format error")
