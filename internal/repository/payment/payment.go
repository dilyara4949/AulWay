package payment

import (
	"aulway/internal/domain"
	"context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, tx *gorm.DB, payment *domain.Payment) error {
	return tx.WithContext(ctx).Create(payment).Error
}

func (r *Repository) UpdateStatus(ctx context.Context, paymentID string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Payment{}).Where("id = ?", paymentID).Update("status", status).Error
}

func (r *Repository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).Where("id = ?", paymentID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *Repository) GetByTicketID(ctx context.Context, ticketID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).
		Where("ticket_id = ?", ticketID).
		First(&payment).Error

	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *Repository) Update(ctx context.Context, tx *gorm.DB, updates map[string]interface{}, paymentID string) error {
	err := tx.WithContext(ctx).
		Model(&domain.Payment{}).
		Where("id = ?", paymentID).
		Updates(updates).Error

	if err != nil {
		return err
	}
	return nil
}
