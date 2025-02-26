package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Email                *string        `gorm:"type:varchar(255);unique" json:"email,omitempty"`
	Phone                *string        `gorm:"type:varchar(20);unique" json:"phone,omitempty"`
	Password             string         `gorm:"type:text;not null" json:"-"`
	FirstName            *string        `gorm:"type:varchar(100)" json:"first_name,omitempty"`
	LastName             *string        `gorm:"type:varchar(100)" json:"last_name,omitempty"`
	RequirePasswordReset bool           `json:"require_password_reset"`
	FirebaseUID          string         `gorm:"type:varchar(128);unique"`
	CreatedAt            time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
	Role                 string         `gorm:"size:20;not null" json:"role"`
}
