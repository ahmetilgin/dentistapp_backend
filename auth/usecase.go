package auth

import (
	"backend/models"
	"context"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUpBusinessUser(ctx context.Context, user* models.BusinessUser) error
	SignUpNormalUser(ctx context.Context, user* models.NormalUser) error
	SignInNormalUser(ctx context.Context, username, password string) (*models.NormalUser, string, error)
	SignInBusinessUser(ctx context.Context, username, password string) (*models.BusinessUser, string, error)
	ParseToken(ctx context.Context, accessToken string) (interface {}, error)
	ResetPasswordNormalUser(ctx context.Context, email string) error	
	ResetPasswordBusinessUser(ctx context.Context, email string) error
}
