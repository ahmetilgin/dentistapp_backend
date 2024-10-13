package country

import (
	"backend/models"
	"context"
)

type Repository interface {
	CreateRegion(ctx context.Context, country *models.Country) error
	Search(ctx context.Context, code, query string) ([]string, error)
}
