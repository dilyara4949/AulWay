package model

import (
	"aulway/internal/handler/user/model"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPassword struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type SigninResponse struct {
	AccessToken string             `json:"access_token"`
	User        model.UserResponse `json:"user"`
}

type SignupResponse struct {
	AccessToken string             `json:"access_token"`
	User        model.UserResponse `json:"user"`
}
