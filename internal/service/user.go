package service

import (
	"aulway/internal/domain"
	"aulway/internal/handler/user/model"
	"aulway/internal/repository/user"
	"context"
	"fmt"
	"github.com/google/uuid"
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
		ID:                   userid,
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
	if req.Email != nil {
		if err := req.ValidateEmail(*req.Email); err != nil {
			return err
		}

		updates["email"] = *req.Email

	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	if len(updates) == 0 {
		return nil
	}

	return service.repo.Update(ctx, updates, id)
}

//func (service *User) CreateUser(ctx context.Context, signup model.SignupRq, password string) (*domain.User, error) {
//	encryptedPassword, err := bcrypt.GenerateFromPassword(
//		[]byte(password),
//		bcrypt.DefaultCost,
//	)
//	if err != nil {
//		return nil, fmt.Errorf("generate password error: %w", err)
//	}
//
//	userid, err := uuid.NewV7()
//	if err != nil {
//		return nil, fmt.Errorf("generate uuid error: %w", err)
//	}
//
//	usr := domain.User{
//		ID:                   userid,
//		Phone:                signup.Phone,
//		Role:                 userRole,
//		Password:             string(encryptedPassword),
//		RequirePasswordReset: false,
//	}
//
//	err = service.repo.Create(ctx, &usr)
//	return &usr, err
//}
