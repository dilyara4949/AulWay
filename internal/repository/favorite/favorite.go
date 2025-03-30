package favorite

import (
	"aulway/internal/domain"
	"aulway/internal/repository/errs"
	uerror "aulway/internal/utils/errs"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (repo *Repository) Add(ctx context.Context, favorite *domain.FavoriteRoute) error {
	if err := repo.db.WithContext(ctx).Create(&favorite).Error; err != nil {
		return fmt.Errorf("add favorite route error: %w", err)
	}
	return nil
}

func (repo *Repository) Remove(ctx context.Context, favorite *domain.FavoriteRoute) error {
	res := repo.db.WithContext(ctx).
		Where("user_id = ? AND route_id = ?", favorite.UserID, favorite.RouteID).
		Delete(&domain.FavoriteRoute{})

	if res.Error != nil {
		return uerror.Err{ErrDesc: "remove favorite route error: %w", Err: res.Error.Error()}
	}

	if res.RowsAffected == 0 {
		return uerror.Err{ErrDesc: "favorite route not found: %w", Err: errs.ErrRecordNotFound.Error()}
	}

	return nil
}

func (repo *Repository) GetByUser(ctx context.Context, userID string) ([]domain.FavoriteRoute, error) {
	var favorites []domain.FavoriteRoute

	if err := repo.db.WithContext(ctx).
		Table("favorite_routes").
		Joins("JOIN routes ON favorite_routes.route_id = routes.id").
		Where("favorite_routes.user_id = ?", userID).
		Order("routes.start_date ASC").
		Find(&favorites).Error; err != nil {
		return nil, uerror.Err{ErrDesc: "get favorites error: %w", Err: err.Error()}
	}

	return favorites, nil
}

func (repo *Repository) Exists(ctx context.Context, userID, routeID string) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&domain.FavoriteRoute{}).
		Where("user_id = ? AND route_id = ?", userID, routeID).
		Count(&count).Error

	if err != nil {
		return false, uerror.Err{ErrDesc: "check favorite exists error: %w", Err: err.Error()}
	}

	return count > 0, nil
}

func (repo *Repository) GetFavoriteRoutes(ctx context.Context, userId string, page, pageSize int) ([]domain.Route, error) {
	routes := make([]domain.Route, 0)

	offset := (page - 1) * pageSize

	err := repo.db.WithContext(ctx).
		Table("favorite_routes").
		Select(`
			routes.id, routes.departure, routes.destination,
			routes.start_date, routes.end_date, routes.available_seats,
			routes.bus_id, routes.price,
			routes.created_at, routes.updated_at
		`).
		Joins("JOIN routes ON favorite_routes.route_id = routes.id").
		Where("favorite_routes.user_id = ?", userId).
		Order("routes.start_date ASC").
		Limit(pageSize).
		Offset(offset).
		Scan(&routes).Error

	if err != nil {
		return nil, err
	}

	return routes, nil

}
