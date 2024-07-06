package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Region represents a region
type Region struct {
    ID     primitive.ObjectID `bson:"_id,omitempty"`
    Name   string             `bson:"name"`
    Code   string             `bson:"code"`
    Cities []City             `bson:"cities"`
}

type City struct {
    ID         primitive.ObjectID `bson:"_id,omitempty"`
    Name       string             `bson:"name"`
    RegionID    primitive.ObjectID `bson:"region_id"`
    Districts  []District         `bson:"districts"`
}

// District represents a district
type District struct {
    ID     primitive.ObjectID `bson:"_id,omitempty"`
    Name   string             `bson:"name"`
    CityID primitive.ObjectID `bson:"city_id"`
}