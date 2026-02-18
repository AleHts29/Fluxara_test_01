package db

import (
	"context"
	"fluxara/internal/domain"
	"fluxara/internal/ports/repos"
)

type DbService struct {
	repo repos.DbReporer
}

func NewDbService(repo repos.DbReporer) *DbService {
	return &DbService{
		repo: repo,
	}
}

// arte
func (db *DbService) GetFullData(ctx context.Context) ([]domain.CareerFull, error) {
	fullData, err := db.repo.GetFullData(ctx)
	if err != nil {
		return fullData, err
	}

	return fullData, err
}
func (db *DbService) GetCarrerasAll(ctx context.Context) ([]domain.CareerFull, error) {
	carrers, err := db.repo.GetCarrerasAll(ctx)
	if err != nil {
		return carrers, err
	}

	return carrers, err
}
