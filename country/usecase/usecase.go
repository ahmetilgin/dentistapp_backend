package usecase

import (
	"backend/country"
	"backend/models"
	"context"
)

type RegionUseCase struct {
	regionRepo country.Repository
}

func NewRegionUseCase(regionRepo country.Repository) *RegionUseCase {
	return &RegionUseCase{
		regionRepo: regionRepo,
	}
}

func (b RegionUseCase) CreateRegion(ctx context.Context, country *models.Country) error {
	return b.regionRepo.CreateRegion(ctx, country)
}

func (b RegionUseCase) Search(ctx context.Context, query, code string) ([]string, error) {
	return b.regionRepo.Search(ctx, query, code)
}
