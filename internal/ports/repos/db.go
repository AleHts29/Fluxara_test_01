package repos

import (
	"context"
	"fluxara/internal/domain"
)

type DbReporer interface {
	GetProduct(ctx context.Context, id string) (domain.Product, error)
	GetProductsAll(ctx context.Context) ([]domain.Product, error)
}
