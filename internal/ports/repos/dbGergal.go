package repos

import (
	"context"
	"fluxara/internal/domain"
)

type DbReporerGergal interface {
	GetCatalog(ctx context.Context) ([]domain.Product, error)
	GetDeliveryZones(ctx context.Context) ([]domain.DeliveryZone, error)
}
