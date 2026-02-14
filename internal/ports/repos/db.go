package repos

import (
	"context"
	"fluxara/internal/domain"
)

type DbReporer interface {
	// GetProductsAll(ctx context.Context) ([]domain.Product, error)
	// GetProduct(ctx context.Context, id string) (domain.Product, error)
	// arte
	GetCarrerasAll(ctx context.Context) ([]domain.Career, error)
	GetCarrerasResumen(ctx context.Context) ([]domain.CareersResumen, error)
	GetCarrerasByName(ctx context.Context, name string) (domain.Career, error)
}
