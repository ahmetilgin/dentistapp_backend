package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// BaseUser model
type BaseUser struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Email    string             `bson:"email" json:"email"`
    Username string             `bson:"username" json:"username"`
    Password string             `bson:"password" json:"password"`
}

// NormalUser model
type NormalUser struct {
    BaseUser
    // Normal kullanıcıya özgü alanlar
}

// BusinessUser model
type BusinessUser struct {
    BaseUser
    BusinessName    string `bson:"business_name" json:"business_name"`
    BusinessAddress string `bson:"business_address" json:"business_address"`
}