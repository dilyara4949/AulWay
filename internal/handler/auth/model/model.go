package model

import "aulway/internal/domain"

type SignupRq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type SigninRq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPassword struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type SignupRs struct {
	AccessToken string      `json:"access_token"`
	User        domain.User `json:"user"`
}
