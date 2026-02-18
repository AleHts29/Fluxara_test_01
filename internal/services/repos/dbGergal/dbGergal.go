package db

import (
	"context"
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
