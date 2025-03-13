package ticket

import (
	"aulway/internal/domain"
	"aulway/internal/repository/errs"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (repo *Repository) BeginTransaction() *gorm.DB {
	return repo.db.Begin()
}

func (repo *Repository) Create(ctx context.Context, tx *gorm.DB, ticket *domain.Ticket) error {
	if err := tx.WithContext(ctx).Create(&ticket).Error; err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Get(ctx context.Context, id string) (*domain.Ticket, error) {
	ticket := new(domain.Ticket)

	if err := repo.db.WithContext(ctx).First(&ticket, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, fmt.Errorf("get bus error: %w", err)
	}

	return ticket, nil
}

func (repo *Repository) Update(ctx context.Context, tx *gorm.DB, updates map[string]interface{}, id string) error {
	err := tx.WithContext(ctx).Model(&domain.Ticket{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, id string) error {
	if err := repo.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Route{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrRecordNotFound
		}
		return fmt.Errorf("delete route error: %w", err)
	}
	return nil
}

func (repo *Repository) GetUpcomingTickets(ctx context.Context, userID string, now time.Time) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	err := repo.db.WithContext(ctx).Where("user_id = ? AND created_at > ?", userID, now).Find(&tickets).Error
	return tickets, err
}

func (repo *Repository) GetPastTickets(ctx context.Context, userID string, now time.Time) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	err := repo.db.WithContext(ctx).Where("user_id = ? AND created_at <= ?", userID, now).Find(&tickets).Error
	return tickets, err
}
