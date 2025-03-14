package service

import (
	"aulway/internal/domain"
	"aulway/internal/handler/auth/model"
	"aulway/internal/repository/user"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"user_role"`
	jwt.RegisteredClaims
}

type Auth struct {
	repo  user.Repository
	redis *redis.Client
}

func (s *Auth) VerifyResetCode(ctx context.Context, req model.VerifyResetCodeRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewAuthService(userRepo user.Repository, redis *redis.Client) *Auth {
	return &Auth{
		repo:  userRepo,
		redis: redis,
	}
}

func (s *Auth) SendResetCode(ctx context.Context, email string) error {
	usr, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	err = s.redis.Set(ctx, "reset_code:"+email, code, 10*time.Minute).Err()
	if err != nil {
		return errors.New("failed to store reset code")
	}

	message := fmt.Sprintf("Your password reset code is: %s. It will expire in 10 minutes.", code)
	return SendEmail(usr.Email, "Password Reset Code", message)
}

//func (s *Auth) VerifyResetCode(ctx context.Context, req model.VerifyResetCodeRequest) error {
//	storedCode, err := s.redis.Get(ctx, "reset_code:"+req.Email).Result()
//	if err == redis.Nil {
//		return errors.New("reset code expired or invalid")
//	} else if err != nil {
//		return err
//	}
//
//	if storedCode != req.Code {
//		return errors.New("invalid reset code")
//	}
//
//	err = s.repo.UpdatePassword(ctx, req.Email, req.NewPassword, false)
//	if err != nil {
//		return errors.New("failed to reset password")
//	}
//
//	s.redis.Del(ctx, "reset_code:"+req.Email)
//
//	return nil
//}

func (service *Auth) CreateAccessToken(ctx context.Context, user domain.User, jwtSecret string, expiry int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expiry) * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (service *Auth) ValidatePhone(phone string) error {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	re := regexp.MustCompile(`^(?:\+7|8)\d{10}$`)

	if !re.MatchString(phone) {
		return errors.New("invalid phone number format")
	}

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
	re := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$`)

	if !re.MatchString(password) {
		return errors.New("password must be at least 8 characters long and include 1 uppercase, 1 lowercase, 1 number, and 1 special character")
	}

	return nil
}
