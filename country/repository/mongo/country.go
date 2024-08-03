package mongo

import (
	"backend/models"
	"context"
	"regexp"
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

func (r *RegionRepository) Search(ctx context.Context, code, query string) ([]string, error) {
	code = strings.ToUpper(code)
	escapedQuery := regexp.QuoteMeta(query)
	if code == "EN" {
		code = "SQ"
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"code": code,
				"$or": []bson.M{
					{"name": bson.M{"$regex": "^" + escapedQuery, "$options": "i"}},
					{"cities.name": bson.M{"$regex": "^" + escapedQuery, "$options": "i"}},
					{"cities.districts.name": bson.M{"$regex": "^" + escapedQuery, "$options": "i"}},
				},
			},
		},
		{
			"$project": bson.M{
				"matches": bson.M{
					"$concatArrays": []interface{}{
						bson.M{"$filter": bson.M{
							"input": []string{"$name"},
							"as":    "name",
							"cond":  bson.M{"$regexMatch": bson.M{"input": "$$name", "regex": "^" + escapedQuery, "options": "i"}},
						}},
						bson.M{"$filter": bson.M{
							"input": "$cities.name",
							"as":    "cityName",
							"cond":  bson.M{"$regexMatch": bson.M{"input": "$$cityName", "regex": "^" + escapedQuery, "options": "i"}},
						}},
						bson.M{"$reduce": bson.M{
							"input":        "$cities",
							"initialValue": []string{},
							"in": bson.M{
								"$concatArrays": []interface{}{
									"$$value",
									bson.M{"$filter": bson.M{
										"input": "$$this.districts.name",
										"as":    "districtName",
										"cond":  bson.M{"$regexMatch": bson.M{"input": "$$districtName", "regex": "^" + escapedQuery, "options": "i"}},
									}},
								},
							},
						}},
					},
				},
			},
		},
		{
			"$unwind": "$matches",
		},
		{
			"$group": bson.M{
				"_id":     nil,
				"results": bson.M{"$addToSet": "$matches"},
			},
		},
		{
			"$project": bson.M{
				"results": bson.M{"$slice": []interface{}{"$results", 10}},
				"_id":     0,
			},
		},
	}

	cursor, err := r.regionCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Results []string `bson:"results"`
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []string{}, nil
	}

	return results[0].Results, nil
}
