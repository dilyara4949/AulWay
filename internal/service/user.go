package service

import (
	"aulway/internal/domain"
	auth "aulway/internal/handler/auth/model"
	"aulway/internal/handler/user/model"
	"aulway/internal/repository/user"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const userRole = "user"

type User struct {
	repo user.Repository
}

func NewUserService(userRepo user.Repository) *User {
	return &User{
		repo: userRepo,
	}
}

func (service *User) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	return service.repo.Get(ctx, id)
}

func (service *User) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return service.repo.GetByEmail(ctx, email)
}

func (service *User) GetUserByFbUid(ctx context.Context, uid string) (*domain.User, error) {
	return service.repo.GetUserByFbUid(ctx, uid)
}

func (service *User) CreateUser(ctx context.Context, email string, uid string) (*domain.User, error) {
	userid, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate uuid error: %w", err)
	}

	usr := domain.User{
		ID:                   userid.String(),
		Email:                &email,
		Role:                 userRole,
		FirebaseUID:          uid,
		RequirePasswordReset: false,
	}

	err = service.repo.Create(ctx, &usr)
	return &usr, err
}

func (service *User) UpdateUser(ctx context.Context, req model.UpdateUserRequest, id string) error {
	updates := make(map[string]interface{})

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	if len(updates) == 0 {
		return nil
	}

	return service.repo.Update(ctx, updates, id)
}

func (service *User) ResetPassword(ctx context.Context, req auth.ResetPassword) error {
	usr, err := service.ValidateUser(ctx, auth.SigninRq{
		Email:    req.Email,
		Password: req.OldPassword,
	})

	if err != nil {
		return err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return fmt.Errorf("generate password error: %w", err)
	}

	err = service.repo.UpdatePassword(ctx, usr.ID, string(encryptedPassword), false)
	return err
}

func (service *User) ValidateUser(ctx context.Context, signin auth.SigninRq) (*domain.User, error) {
	usr, err := service.repo.GetByEmail(ctx, signin.Email)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(signin.Password)) != nil {
		return nil, errors.New("invalid password")
	}

	return usr, nil
}

func (service *User) GetUsers(ctx context.Context, page, pageSize int) ([]domain.User, error) {
	return service.repo.GetUsers(ctx, page, pageSize)
}
