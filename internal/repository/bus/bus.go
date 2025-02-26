package bus

import (
	"aulway/internal/domain"
	"aulway/internal/repository/errs"
	uerror "aulway/internal/utils/errs"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (repo *Repository) Create(ctx context.Context, bus *domain.Bus) error {
	if err := repo.db.WithContext(ctx).Create(&bus).Error; err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Get(ctx context.Context, id string) (*domain.Bus, error) {
	bus := new(domain.Bus)

	if err := repo.db.WithContext(ctx).First(&bus, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, fmt.Errorf("get bus error: %w", err)
	}

	return bus, nil
}

func (repo *Repository) GetByNumber(ctx context.Context, number string) (*domain.Bus, error) {
	bus := new(domain.Bus)

	if err := repo.db.WithContext(ctx).First(bus, "number = ?", number).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, err
	}

	return bus, nil
}

func (repo *Repository) Update(ctx context.Context, bus *domain.Bus) error {
	if err := repo.db.WithContext(ctx).Save(bus).Error; err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res := repo.db.WithContext(ctx).Delete(&domain.Bus{}, id)
	if res.Error != nil {
		return uerror.Err{ErrDesc: "delete bus error: %w", Err: res.Error.Error()}
	}

	if res.RowsAffected == 0 {
		return uerror.Err{ErrDesc: "delete bus error: %w", Err: errs.ErrRecordNotFound.Error()}
	}

	return nil
}

func (repo *Repository) GetBuses(ctx context.Context, page, pageSize int) ([]domain.Bus, error) {
	var buses []domain.Bus

	offset := (page - 1) * pageSize

	if err := repo.db.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&buses).Error; err != nil {
		return nil, uerror.Err{ErrDesc: "get page error: %w", Err: err.Error()}
	}

	return buses, nil
}
