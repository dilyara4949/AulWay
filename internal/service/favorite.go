package service

import (
	"aulway/internal/domain"
	"aulway/internal/repository/favorite"
	"context"
)

type Favorite struct {
	repo favorite.Repository
}

func NewFavoriteService(favoriteRepo favorite.Repository) *Favorite {
	return &Favorite{
		repo: favoriteRepo,
	}
}

func (s *Favorite) Add(ctx context.Context, favorite *domain.FavoriteRoute) error {
	return s.repo.Add(ctx, favorite)
}

func (s *Favorite) Remove(ctx context.Context, favorite *domain.FavoriteRoute) error {
	return s.repo.Remove(ctx, favorite)
}

func (s *Favorite) GetFavoriteRoutes(ctx context.Context, userId string, page, pageSize int) ([]domain.Route, error) {
	return s.repo.GetFavoriteRoutes(ctx, userId, page, pageSize)
}
