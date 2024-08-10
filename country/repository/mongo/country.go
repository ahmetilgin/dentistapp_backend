package mongo

import (
	"backend/models"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegionRepository struct {
	db        *mongo.Database
	regionCol *mongo.Collection
}

func NewRegionRepository(db *mongo.Database, countriesCollection string) *RegionRepository {
	return &RegionRepository{
		db:        db,
		regionCol: db.Collection(countriesCollection),
	}
}

func (r *RegionRepository) CreateRegion(ctx context.Context, country *models.Country) error {
	_, err := r.regionCol.InsertOne(ctx, country)
	return err
}

func (r *RegionRepository) Search(ctx context.Context, query, code string) ([]string, error) {
	code = strings.ToUpper(code)
	if code == "EN" {
		code = "SQ"
	}

	pipeline := []bson.M{
		{"$match": bson.M{"code": code}},
		{"$unwind": "$cities"},
		{"$unwind": "$cities.districts"},
		{"$or": []bson.M{
			{"$match": bson.M{"cities.name": bson.M{"$regex": "^" + query, "$options": "i"}}},
			{"$match": bson.M{"cities.districts.name": bson.M{"$regex": "^" + query, "$options": "i"}}},
		}},
		{"$project": []bson.M{
			{"city_name": "$cities.name"},
			{"district_name": "$cities.districts.name"},
		}},
	}

	cursor, err := r.regionCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		CityName string `bson:"city_name"`
		District string `bson:"district_name"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	stringResults := make([]string, len(results))
	for i, result := range results {
		stringResults[i] = result.CityName
	}

	return stringResults, nil
}
