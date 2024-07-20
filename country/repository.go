package country

import (
	"backend/models"
	"context"
)

type Repository interface {
	CreateRegion(ctx context.Context, country *models.Country) error
	Search(ctx context.Context, query, code string) ([]string, error)
}
