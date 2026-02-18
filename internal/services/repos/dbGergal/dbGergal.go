package db

import (
	"context"
	"errors"
	"fluxara/internal/domain"
	"fluxara/internal/ports/repos"
)

type DbService struct {
	repo repos.DbReporerGergal
}

func NewDbServiceGergal(repo repos.DbReporerGergal) *DbService {
	return &DbService{
		repo: repo,
	}
}

func (db *DbService) GetCatalog(ctx context.Context) ([]domain.Product, error) {
	fullCatalog, err := db.repo.GetCatalog(ctx)
	if err != nil {
		return fullCatalog, err
	}

	return fullCatalog, err
}

func (db *DbService) GetDeliveryZones(ctx context.Context) ([]domain.DeliveryZone, error) {
	deliveryZones, err := db.repo.GetDeliveryZones(ctx)
	if err != nil {
		return deliveryZones, err
	}

	return deliveryZones, err
}

func (db *DbService) CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("la orden debe tener al menos un item")
	}

	for _, it := range req.Items {
		if it.Quantity <= 0 {
			return nil, errors.New("cantidad inválida")
		}
	}

	return db.repo.CreateOrder(ctx, req)
}
