package service

import (
	"aulway/internal/domain"
	"aulway/internal/handler/auth/model"
	"aulway/internal/repository/user"
	"aulway/internal/utils/config"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
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
	smpt  config.SMTP
}

func NewAuthService(userRepo user.Repository, redis *redis.Client, smtp config.SMTP) *Auth {
	return &Auth{
		repo:  userRepo,
		redis: redis,
		smpt:  smtp,
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

	message := fmt.Sprintf("Aulway\r\nYour password reset code is: %s. It will expire in 10 minutes.", code)
	return SendEmail(usr.Email, "Password Reset Code", message, s.smpt)
}

func (s *Auth) VerifyResetCode(ctx context.Context, req model.VerifyResetCodeRequest) error {
	storedCode, err := s.redis.Get(ctx, "reset_code:"+req.Email).Result()
	if err == redis.Nil {
		return errors.New("reset code invalid or email not found")
	} else if err != nil {
		return err
	}

	if storedCode != req.Code {
		return errors.New("invalid reset code")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)

	err = s.repo.UpdatePassword(ctx, req.Email, string(encryptedPassword), false)
	if err != nil {
		return errors.New("failed to reset password")
	}

	s.redis.Del(ctx, "reset_code:"+req.Email)

	return nil
}

func (s *Auth) CreateAccessToken(ctx context.Context, user domain.User, jwtSecret string, expiry int) (string, error) {
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

func (s *Auth) ValidatePhone(phone string) error {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	re := regexp.MustCompile(`^(?:\+7|8)\d{10}$`)

	if !re.MatchString(phone) {
		return errors.New("invalid phone number format")
	}

	return nil
}
