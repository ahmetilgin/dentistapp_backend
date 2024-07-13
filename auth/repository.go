package auth

import (
	"backend/models"
	"context"
)

type UserRepository interface {
	CreateBusinessUser(ctx context.Context, user *models.BusinessUser) error
	CreateNormalUser(ctx context.Context, user *models.NormalUser) error
	GetBusinessUser(ctx context.Context, username, password string) (*models.BusinessUser, error)
	GetNormalUser(ctx context.Context, username, password string) (*models.NormalUser, error)
}
	