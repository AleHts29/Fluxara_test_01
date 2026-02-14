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

// func (db *DbService) GetProductsAll(ctx context.Context) ([]domain.Product, error) {
// 	device, err := db.repo.GetProductsAll(ctx)
// 	if err != nil {
// 		return device, err
// 	}

// 	return device, err
// }

// func (db *DbService) GetProduct(ctx context.Context, id string) (domain.Product, error) {
// 	device, err := db.repo.GetProduct(ctx, id)
// 	if err != nil {
// 		return device, err
// 	}

// 	return device, err
// }

// arte
func (db *DbService) GetCarrerasAll(ctx context.Context) ([]domain.Career, error) {
	carrers, err := db.repo.GetCarrerasAll(ctx)
	if err != nil {
		return carrers, err
	}

	return carrers, err
}

func (db *DbService) GetCarrerasResumen(ctx context.Context) ([]domain.CareersResumen, error) {
	carrers, err := db.repo.GetCarrerasResumen(ctx)
	if err != nil {
		return carrers, err
	}

	return carrers, err
}

func (db *DbService) GetCarrerasByName(ctx context.Context, name string) (domain.Career, error) {
	career, err := db.repo.GetCarrerasByName(ctx, name)
	if err != nil {
		return career, err
	}

	return career, err
}
