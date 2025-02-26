package service

import (
	"aulway/internal/repository/user"
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type Auth struct {
	repo user.Repository
}

func NewAuthService(userRepo user.Repository) *Auth {
	return &Auth{
		repo: userRepo,
	}
}

func (service *Auth) VerifyFirebaseToken(client *auth.Client, idToken string) (*auth.Token, error) {
	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

//func (service *Auth) CreateAccessToken(ctx context.Context, user domain.User, jwtSecret string, expiry int) (string, error) {
//	expirationTime := time.Now().Add(time.Duration(expiry) * time.Hour)
//	claims := &Claims{
//		UserID: user.ID,
//		Role:   user.Role,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(expirationTime),
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
//	accessToken, err := token.SignedString([]byte(jwtSecret))
//	if err != nil {
//		return "", err
//	}
//
//	return accessToken, nil
//}

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
