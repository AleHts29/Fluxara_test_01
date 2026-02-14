package repos

import (
	"context"
	"fluxara/internal/domain"
)

type DbReporer interface {
	// GetProductsAll(ctx context.Context) ([]domain.Product, error)
	// GetProduct(ctx context.Context, id string) (domain.Product, error)
	// arte
	GetFullData(ctx context.Context) ([]domain.CareerFull, error)
	GetCarrerasAll(ctx context.Context) ([]domain.CareerFull, error)
	GetCarrerasResumen(ctx context.Context) ([]domain.CareerFull, error)
}
