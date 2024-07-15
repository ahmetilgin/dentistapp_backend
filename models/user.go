package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasswordResetToken struct {
    Token     string             `bson:"token"`
    UserID    primitive.ObjectID `bson:"user_id"`
    ExpiresAt time.Time          `bson:"expires_at"`
}

// BaseUser model
type NormalUser struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Email    string             `bson:"email" json:"email"`
    Username string             `bson:"username" json:"username"`
    Password string             `bson:"password" json:"password"`
}

// BusinessUser model
type BusinessUser struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Email    string             `bson:"email" json:"email"`
    Username string             `bson:"username" json:"username"`
    Password string             `bson:"password" json:"password"`
    BusinessName    string      `bson:"businessName" json:"businessName"`
    BusinessAddress string      `bson:"businessAddress" json:"businessAddress"`
}