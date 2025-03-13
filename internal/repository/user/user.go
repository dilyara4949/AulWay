package user

import (
	"aulway/internal/domain"
	"aulway/internal/repository/errs"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (repo *Repository) Create(ctx context.Context, user *domain.User) error {
	if err := repo.db.WithContext(ctx).Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			slog.Debug(err.Error())
			return errs.EmailAlreadyExists
		}
		return fmt.Errorf("create user error: %w", err)
	}

	return nil
}

func (repo *Repository) Get(ctx context.Context, id string) (*domain.User, error) {
	var user *domain.User

	if err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, fmt.Errorf("get user error: %w", err)
	}

	return user, nil
}

func (repo *Repository) Update(ctx context.Context, updates map[string]interface{}, id string) error {
	err := repo.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, id string) error {
	now := time.Now()
	res := repo.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("deleted_at", now)

	if res.Error != nil {
		return fmt.Errorf("soft delete user error: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return errs.ErrRecordNotFound
	}

	return nil
}

func (repo *Repository) GetUsers(ctx context.Context, page, pageSize int) ([]domain.User, error) {
	var users []domain.User

	offset := (page - 1) * pageSize

	if err := repo.db.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("get all users error: %w", err)
	}

	return users, nil
}

func (repo *Repository) UpdatePassword(ctx context.Context, userID string, newPassword string, requirePasswordReset bool) error {
	return repo.db.WithContext(ctx).
		Model(&domain.User{ID: userID}).
		Updates(map[string]interface{}{
			"password":               newPassword,
			"require_password_reset": requirePasswordReset,
		}).Error
}

func (repo *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	if err := repo.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, fmt.Errorf("get user by email error: %w", err)
	}

	return &user, nil
}

func (repo *Repository) GetUserByFbUid(ctx context.Context, uid string) (*domain.User, error) {
	var user domain.User

	if err := repo.db.WithContext(ctx).First(&user, "firebase_uid = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, fmt.Errorf("get user by fbUid error: %w", err)
	}

	return &user, nil
}
