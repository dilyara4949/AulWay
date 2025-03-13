package domain

import "time"

type Page struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"unique;not null"` // "about_us", "privacy_policy", "support"
	Content   string    `gorm:"not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
