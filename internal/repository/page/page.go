package page

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

func (repo *Repository) GetPage(ctx context.Context, title string) (domain.Page, error) {
	var page domain.Page
	if err := repo.db.WithContext(ctx).Where("title = ?", title).First(&page).Error; err != nil {
		return page, err
	}
	return page, nil
}

func (repo *Repository) UpdatePage(ctx context.Context, title, content string) error {
	return repo.db.WithContext(ctx).Model(&domain.Page{}).
		Where("title = ?", title).
		Update("content", content).Error
}
