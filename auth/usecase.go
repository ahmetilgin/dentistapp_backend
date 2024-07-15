package auth

import (
	"backend/models"
	"context"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUpBusinessUser(ctx context.Context, user *models.BusinessUser) error
	SignUpNormalUser(ctx context.Context, user *models.NormalUser) error
	SignInNormalUser(ctx context.Context, username, password string) (*models.NormalUser, string, error)
	SignInBusinessUser(ctx context.Context, username, password string) (*models.BusinessUser, string, error)
	ParseToken(ctx context.Context, accessToken string) (interface{}, error)
	SendEmailNormalUser(ctx context.Context, host, email string) error
	SendEmailBusinessUser(ctx context.Context, host, email string) error
	ResetPassword(ctx context.Context, user interface{}, token, newPassword string) error
}
