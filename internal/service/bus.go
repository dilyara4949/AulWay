package service

import (
	"aulway/internal/domain"
	"aulway/internal/handler/bus/model"
	"aulway/internal/repository/bus"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Bus struct {
	repo bus.Repository
}

func NewBusService(busRepo bus.Repository) *Bus {
	return &Bus{
		repo: busRepo,
	}
}

func (service *Bus) CreateBus(ctx context.Context, request model.CreateRequest) (*domain.Bus, error) {
	busId, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate uuid error: %w", err)
	}

	response := &domain.Bus{
		Id:         busId.String(),
		Number:     request.Number,
		TotalSeats: request.TotalSeats,
	}

	err = service.repo.Create(ctx, response)
	return response, err
}

func (service *Bus) GetByNumber(ctx context.Context, number string) (*domain.Bus, error) {
	return service.repo.GetByNumber(ctx, number)
}

func (service *Bus) Get(ctx context.Context, id string) (*domain.Bus, error) {
	return service.repo.Get(ctx, id)
}

func (service *Bus) GetBusesList(ctx context.Context, page, pageSize int) ([]domain.Bus, error) {
	return service.repo.GetBusesList(ctx, page, pageSize)
}
