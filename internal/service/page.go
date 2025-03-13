package service

import (
	"aulway/internal/domain"
	"aulway/internal/repository/page"
	"context"
)

type Page struct {
	repo page.Repository
}

func NewPageService(pageRepo page.Repository) *Page {
	return &Page{
		repo: pageRepo,
	}
}

func (s *Page) GetPage(ctx context.Context, title string) (domain.Page, error) {
	return s.repo.GetPage(ctx, title)
}

func (s *Page) UpdatePage(ctx context.Context, title, content string) error {
	return s.repo.UpdatePage(ctx, title, content)
}
