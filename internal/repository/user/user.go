package user

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if err := repo.db.WithContext(ctx).Create(&user).Error; err != nil {
		return fmt.Errorf("create user error: %w", err)
	}

	return nil
}

func (repo *UserRepository) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user domain.User

	if err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errs.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf("get user error: %w", err)
	}

	return user, nil
}

func (repo *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	if err := repo.db.WithContext(ctx).Save(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("update user error: %w", err)
	}

	return user, nil
}

func (repo *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := repo.db.WithContext(ctx).Delete(&domain.User{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete user error: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (repo *UserRepository) GetUsers(ctx context.Context, page, pageSize int) ([]domain.User, error) {
	var users []domain.User

	offset := (page - 1) * pageSize

	if err := repo.db.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("get all users error: %w", err)
	}

	return users, nil
}

func (repo *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string, requirePasswordReset bool) error {
	return repo.db.WithContext(ctx).
		Model(&domain.User{ID: userID}).
		Updates(map[string]interface{}{
			"password":               newPassword,
			"require_password_reset": requirePasswordReset,
		}).Error
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	if err := repo.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by email error: %w", err)
	}

	return &user, nil
}
