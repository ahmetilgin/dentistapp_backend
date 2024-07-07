package usecase

import (
	"backend/models"
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"backend/auth"

	"github.com/dgrijalva/jwt-go/v4"
)

type AuthClaims struct {
	jwt.StandardClaims
	BusinessUser *models.BusinessUser `json:"business_user"`
	NormalUser *models.NormalUser `json:"normal_user"`

}

type AuthUseCase struct {
	userRepo       		auth.UserRepository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(
	userRepo auth.UserRepository,
	hashSalt string,
	signingKey []byte,
	tokenTTLSeconds time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
	}
}

func (a *AuthUseCase) SignUpBusinessUser(ctx context.Context, user *models.BusinessUser) error {
	pwd := sha1.New()
	
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil));

	return a.userRepo.CreateBusinessUser(ctx, user)
}

func (a *AuthUseCase) SignUpNormalUser(ctx context.Context, user *models.NormalUser) error {
	pwd := sha1.New()
	
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil));

	return a.userRepo.CreateNormalUser(ctx, user)
}

func (a *AuthUseCase) SignInNormalUser(ctx context.Context, username, password string) (*models.NormalUser, string, error)  {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userRepo.GetNormalUser(ctx, username, password)
	
	if err != nil {
		return nil,"", auth.ErrUserNotFound
	}

	claims := AuthClaims{
		NormalUser: user,
		BusinessUser: nil,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.signingKey)
    if err != nil {
        return nil, "", err
    }
	return user, signedToken, nil
}

func (a *AuthUseCase) SignInBusinessUser(ctx context.Context, username, password string) (*models.BusinessUser, string, error)  {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userRepo.GetBusinessUser(ctx, username, password)
	
	if err != nil {
		return nil,"", auth.ErrUserNotFound
	}

	claims := AuthClaims{
		NormalUser: nil,
		BusinessUser: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.signingKey)
    if err != nil {
        return nil, "", err
    }
	return user, signedToken, nil
}




func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (interface {}, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		if (claims.BusinessUser == nil && claims.NormalUser == nil) {
			return nil, auth.ErrInvalidAccessToken
		}
		if (claims.BusinessUser == nil) {
			return claims.NormalUser, nil
		}

		return claims.BusinessUser, nil
	}

	return nil, auth.ErrInvalidAccessToken
}
