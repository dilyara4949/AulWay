package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                   string         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Email                string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	Phone                string         `gorm:"type:varchar(20);unique;not null" json:"phone"`
	Password             string         `gorm:"type:text;not null" json:"-"`
	FirstName            string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName             string         `gorm:"type:varchar(100);not null" json:"last_name"`
	RequirePasswordReset bool           `gorm:"default:false" json:"require_password_reset"`
	CreatedAt            time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
	Role                 string         `gorm:"size:20;not null" json:"role"`
}
