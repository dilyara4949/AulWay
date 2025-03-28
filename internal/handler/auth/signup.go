package auth

import (
	"aulway/internal/domain"
	"aulway/internal/handler/auth/model"
	"aulway/internal/handler/user"
	usermodel "aulway/internal/handler/user/model"
	"aulway/internal/service"
	"aulway/internal/utils/config"
	uerrs "aulway/internal/utils/errs"
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type Service interface {
	CreateAccessToken(ctx context.Context, user domain.User, jwtSecret string, expiry int) (string, error)
	SendResetCode(ctx context.Context, email string) error
	VerifyResetCode(ctx context.Context, req model.VerifyResetCodeRequest) error
}

// ForgotPasswordHandler
// @Summary Forgot password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.ForgotPasswordRequest true "forgot password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errs.Err "Bad Request - Invalid request body"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /auth/forgot-password [post]
func ForgotPasswordHandler(svc Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.ForgotPasswordRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "Invalid request", ErrDesc: err.Error()})
		}

		err := svc.SendResetCode(c.Request().Context(), req.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Reset code sent to email"})
	}
}

// VerifyForgotPasswordHandler
// @Summary verify forgot password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.VerifyResetCodeRequest true "required"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errs.Err "Bad Request - Invalid request body"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /auth/forgot-password/verify [post]
func VerifyForgotPasswordHandler(svc Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.VerifyResetCodeRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "Invalid request"})
		}

		err := svc.VerifyResetCode(c.Request().Context(), req)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, uerrs.Err{Err: err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successful"})
	}
}

// SignupHandler
// @Summary User Signup
// @Description Register a new user and send verification code.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.SignupRequest true "Signup Request Body"
// @Success 200 {string} string "response"
// @Failure 400 {object} errs.Err "Bad Request - Invalid request body"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /auth/signup [post]
func SignupHandler(redisClient *redis.Client, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.SignupRequest

		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "incorrect req body", ErrDesc: err.Error()})
		}

		if req.Password == "" || req.Email == "" {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "signup failed", ErrDesc: "fields cannot be empty"})
		}

		if err = ValidateEmail(req.Email); err != nil {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "signup failed", ErrDesc: err.Error()})
		}

		if err = ValidatePassword(req.Password); err != nil {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "signup failed", ErrDesc: err.Error()})
		}

		verificationCode := fmt.Sprintf("%06d", rand.Intn(1000000))

		err = redisClient.Set(c.Request().Context(), "email_verification:"+req.Email, verificationCode, 10*time.Minute).Err()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "signup failed", ErrDesc: "failed to store verification code"})
		}

		err = service.SendEmail(req.Email, "Email Verification Code", fmt.Sprintf("Your verification code is: %s", verificationCode), cfg.SMTP)
		if err != nil {
			log.Print(err.Error())
			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "signup failed", ErrDesc: "failed to send verification email"})
		}

		return c.JSON(http.StatusOK, echo.Map{"message": "Verification code sent to email. Please verify your email before logging in."})
	}
}

// VerifyEmailHandler
// @Summary User Signup Verification
// @Description Register a new user and verify code.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.VerifyEmailRequest true "VerifyEmail Request Body"
// @Success 200 {object} model.SignupResponse "response"
// @Failure 400 {object} errs.Err "Bad Request - Invalid request body"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /auth/signup/verify [post]
func VerifyEmailHandler(redisClient *redis.Client, userService user.Service, authService Service, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.VerifyEmailRequest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "verification failed", ErrDesc: "invalid request"})
		}

		storedCode, err := redisClient.Get(c.Request().Context(), "email_verification:"+req.Email).Result()
		if err == redis.Nil {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "verification failed", ErrDesc: "verification code expired or not found"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "verification failed", ErrDesc: "server error"})
		}

		if req.Code != storedCode {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "verification failed", ErrDesc: "invalid verification code"})
		}

		createUserModel := model.SignupRequest{
			Email:    req.Email,
			Password: req.Password,
		}

		usr, err := userService.CreateUser(c.Request().Context(), createUserModel)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "verification failed", ErrDesc: "failed to update user status"})
		}

		token, err := authService.CreateAccessToken(c.Request().Context(), *usr, cfg.JWTTokenSecret, cfg.AccessTokenExpire)
		if err != nil {
			slog.Error("signup: error at creating access token,", "error", err.Error())

			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "signup failed", ErrDesc: err.Error()})
		}

		resp := model.SignupResponse{
			AccessToken: token,
			User:        domainUserToResponse(*usr),
		}

		redisClient.Del(c.Request().Context(), "email_verification:"+req.Email)

		return c.JSON(http.StatusOK, resp)
	}
}

//func SignupHandler(authService Service, userService user.Service, cfg config.Config) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		var req model.SignupRequest
//
//		err := c.Bind(&req)
//		if err != nil {
//			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "incorrect req body", ErrDesc: err.Error()})
//		}
//
//		if req.Password == "" || req.Email == "" {
//			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "signup failed", ErrDesc: "fields cannot be empty"})
//		}
//
//		if err = ValidateEmail(req.Email); err != nil {
//			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "signup failed", ErrDesc: err.Error()})
//		}
//
//		if err = ValidatePassword(req.Password); err != nil {
//			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "signup failed", ErrDesc: err.Error()})
//		}
//
//		usr, err := userService.CreateUser(c.Request().Context(), req)
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "signup failed", ErrDesc: err.Error()})
//		}
//
//		token, err := authService.CreateAccessToken(c.Request().Context(), *usr, cfg.JWTTokenSecret, cfg.AccessTokenExpire)
//		if err != nil {
//			slog.Error("signup: error at creating access token,", "error", err.Error())
//
//			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "signup failed", ErrDesc: err.Error()})
//		}
//
//		resp := model.SignupResponse{
//			AccessToken: token,
//			User:        domainUserToResponse(*usr),
//		}
//
//		return c.JSON(http.StatusOK, resp)
//	}
//}

func domainUserToResponse(user domain.User) usermodel.UserResponse {
	return usermodel.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

//func domainUsersToResponse(users []domain.User) []response.User {
//	res := make([]response.User, 0)
//
//	for _, user := range users {
//		res = append(res, response.User{
//			ID:        user.ID,
//			Email:     user.Email,
//			Phone:     user.Phone,
//			CreatedAt: user.CreatedAt,
//			UpdatedAt: user.UpdatedAt,
//		})
//	}
//	return res
//}

func AssignRole(authClient *auth.Client, uid string, role string) error {
	claims := map[string]interface{}{
		"role": role,
	}
	err := authClient.SetCustomUserClaims(context.Background(), uid, claims)
	if err != nil {
		log.Printf("Error setting role for user %s: %v", uid, err)
		return err
	}
	log.Printf("Role '%s' assigned to user %s", role, uid)
	return nil
}

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
