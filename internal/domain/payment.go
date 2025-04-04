package domain

import "time"

type Payment struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	UserID        string    `json:"user_id" gorm:"not null"`
	Amount        int       `json:"amount" gorm:"not null"`
	Status        string    `json:"status" gorm:"not null"` // pending, successful, failed, refunded
	TransactionID string    `json:"transaction_id" gorm:"unique;not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Payment) TableName() string {
	return "payments"
}
