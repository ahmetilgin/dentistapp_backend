package usecase

import (
	"backend/email_service"
	"backend/models"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"backend/auth"

	"github.com/dgrijalva/jwt-go/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	emailService *email_service.EmailService
}

func NewAuthUseCase(
	userRepo auth.UserRepository,
	hashSalt string,
	signingKey []byte,
	tokenTTLSeconds time.Duration,
	emailService *email_service.EmailService,
	) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
		emailService: emailService,
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

func GenerateResetToken() (string, error) {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func CreatePasswordResetToken(userID primitive.ObjectID) (*models.PasswordResetToken, error) {
    token, err := GenerateResetToken()
    if err != nil {
        return nil, err
    }

    resetToken := &models.PasswordResetToken{
        Token:     token,
        UserID:    userID,
        ExpiresAt: time.Now().Add(1 * time.Hour), // Token 1 saat geçerli
    }

	return resetToken, nil
}


func (a *AuthUseCase) ResetPasswordNormalUser(ctx context.Context, email string) error {
	user, err := a.userRepo.GetNormalUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	resetToken, err := CreatePasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	err = a.userRepo.InsetPasswordResetToken(ctx, resetToken)
	if err != nil {
		return err
	}

	err = a.emailService.SendEmail(email, "Password Reset", "Your password reset token is: " + resetToken.Token)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthUseCase) ResetPasswordBusinessUser(ctx context.Context, email string) error {
	user, err := a.userRepo.GetBusinessUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	resetToken, err := CreatePasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	err = a.userRepo.InsetPasswordResetToken(ctx, resetToken)
	if err != nil {
		return err
	}

	err = a.emailService.SendEmail(email, "Password Reset", "Your password reset token is: " + resetToken.Token)
	if err != nil {
		return err
	}
	
	return nil
}