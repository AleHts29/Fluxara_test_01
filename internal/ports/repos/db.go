package repos

import (
	"context"
	"fluxara/internal/domain"
)

type DbReporer interface {
	GetFullData(ctx context.Context) ([]domain.CareerFull, error)
	GetCarrerasAll(ctx context.Context) ([]domain.CareerFull, error)
}
