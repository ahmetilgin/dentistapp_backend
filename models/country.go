package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Country represents a country
type Country struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name"`
	Code   string             `bson:"code"`
	Cities []City             `bson:"cities"`
}

type City struct {
	Name      string     `bson:"name"`
	Districts []District `bson:"districts"`
	Latitude  float64    `bson:"latitude"`
	Longitude float64    `bson:"longitude"`
}

// District represents a district
type District struct {
	Name      string  ` bson:"name"`
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
}
