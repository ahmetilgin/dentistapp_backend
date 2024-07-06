package usecase

import (
	"backend/models"
	"backend/region"
	"context"
)

type RegionUseCase struct {
	regionRepo region.Repository
}

func NewRegionUseCase(regionRepo region.Repository) *RegionUseCase {
	return &RegionUseCase{
		regionRepo: regionRepo,
	}
}

func (b RegionUseCase) CreateRegion(ctx context.Context, region *models.Region) error {
	return b.regionRepo.CreateRegion(ctx, region)
}

func (b RegionUseCase) CreateCity(ctx context.Context, city *models.City) error {
	return b.regionRepo.CreateCity(ctx, city)
}

func (b RegionUseCase) CreateDistrict(ctx context.Context, district *models.District) error {
	return b.regionRepo.CreateDistrict(ctx, district)
}

func (b RegionUseCase) Search(ctx context.Context, query string)  ([]string, error)  {
	return b.regionRepo.Search(ctx, query)
}


