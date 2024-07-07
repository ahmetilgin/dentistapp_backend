package models

import "go.mongodb.org/mongo-driver/bson/primitive"

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