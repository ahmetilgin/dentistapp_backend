package auth

import (
	"backend/models"
	"context"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUpBusinessUser(ctx context.Context, user* models.BusinessUser) error
	SignUpNormalUser(ctx context.Context, user* models.NormalUser) error
	SignIn(ctx context.Context, username, password string) (interface{}, string, error)
	ParseToken(ctx context.Context, accessToken string) (interface {}, error)
}
