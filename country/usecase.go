package country

import (
	"backend/models"
	"context"
)

type UseCase interface {
	CreateRegion(ctx context.Context, country *models.Country) error
	Search(ctx context.Context, query, code string) ([]string, error)
}
