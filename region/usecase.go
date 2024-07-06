package region

import (
	"backend/models"
	"context"
)

type UseCase interface {
	CreateRegion(ctx context.Context,  region* models.Region) error
	CreateCity(ctx context.Context,  city* models.City) error
	CreateDistrict(ctx context.Context,  district* models.District) error
	Search(ctx context.Context, query string)  ([]string, error) 
}
