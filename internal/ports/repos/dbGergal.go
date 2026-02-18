package repos

import (
	"context"
	"fluxara/internal/domain"
)

type DbReporerGergal interface {
	GetCatalog(ctx context.Context) ([]domain.Product, error)
	GetDeliveryZones(ctx context.Context) ([]domain.DeliveryZone, error)
	MarkOrderPaid(ctx context.Context, orderID int) error
	CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error)
}
