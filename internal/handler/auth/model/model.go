package model

import (
	"aulway/internal/handler/user/model"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
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

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyResetCodeRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required"`
}
