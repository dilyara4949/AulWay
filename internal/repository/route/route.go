package route

import (
	"aulway/internal/domain"
	"aulway/internal/repository/errs"
	uerror "aulway/internal/utils/errs"
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

func (repo *Repository) Create(ctx context.Context, route *domain.Route) error {
	if err := repo.db.WithContext(ctx).Create(&route).Error; err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Update(ctx context.Context, updates map[string]interface{}, id string) error {
	err := repo.db.WithContext(ctx).Model(&domain.Route{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}

	return nil
}

func (repo *Repository) UpdateSeat(ctx context.Context, tx *gorm.DB, updates map[string]interface{}, id string) error {
	err := tx.WithContext(ctx).Model(&domain.Route{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}

	return nil
}

func (repo *Repository) Get(ctx context.Context, id string) (*domain.Route, error) {
	route := new(domain.Route)

	if err := repo.db.WithContext(ctx).First(&route, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, fmt.Errorf("get bus error: %w", err)
	}

	return route, nil
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

func (repo *Repository) GetRoutesList(ctx context.Context, departure, destination string, date time.Time, passengers, page, pageSize int) ([]domain.Route, int, error) {
	routes := make([]domain.Route, 0)

	var total int

	offset := (page - 1) * pageSize

	query := `
		SELECT *, COUNT(*) OVER() AS total_count
		FROM routes
		WHERE departure = ? 
		  AND destination = ? 
		  AND start_date >= ? 
		  AND available_seats >= ?
		ORDER BY start_date ASC
		LIMIT ? OFFSET ?
	`

	rows, err := repo.db.WithContext(ctx).Raw(query, departure, destination, date, passengers, pageSize, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var route domain.Route
		if err := rows.Scan(
			&route.Id, &route.Departure, &route.Destination, &route.StartDate, &route.EndDate,
			&route.AvailableSeats, &route.BusId, &route.Price, &route.CreatedAt, &route.UpdatedAt, &total,
		); err != nil {
			return nil, 0, err
		}
		routes = append(routes, route)
	}

	return routes, total, nil
}

func (repo *Repository) GetAllRoutesList(ctx context.Context, page, pageSize int) ([]domain.Route, error) {

	var routes []domain.Route

	offset := (page - 1) * pageSize

	if err := repo.db.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&routes).Error; err != nil {
		return nil, uerror.Err{ErrDesc: "get page error: %w", Err: err.Error()}
	}

	return routes, nil
}
