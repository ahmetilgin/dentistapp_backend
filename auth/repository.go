package auth

import (
	"backend/models"
	"context"
)

type UserRepository interface {
	CreateBusinessUser(ctx context.Context, user *models.BusinessUser) error
	CreateNormalUser(ctx context.Context, user *models.NormalUser) error
	GetUser(ctx context.Context, username, password string) (interface{}, error)
}
