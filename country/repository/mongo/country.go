package mongo

import (
	"backend/models"
	"context"

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
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"code": code,
				"$or": []bson.M{
					{"name": bson.M{"$regex": query, "$options": "i"}},
					{"cities.name": bson.M{"$regex": query, "$options": "i"}},
					{"cities.districts.name": bson.M{"$regex": query, "$options": "i"}},
				},
			},
		},
		{
			"$project": bson.M{
				"result": bson.M{
					"$concatArrays": []interface{}{
						bson.M{"$cond": []interface{}{
							bson.M{"$regexMatch": bson.M{"input": "$name", "regex": query, "options": "i"}},
							[]string{"$name"},
							[]string{},
						}},
						bson.M{"$reduce": bson.M{
							"input":        "$cities",
							"initialValue": []string{},
							"in": bson.M{
								"$concatArrays": []interface{}{
									"$$value",
									bson.M{"$cond": []interface{}{
										bson.M{"$regexMatch": bson.M{"input": "$$this.name", "regex": query, "options": "i"}},
										[]string{"$$this.name"},
										[]string{},
									}},
									bson.M{"$reduce": bson.M{
										"input":        "$$this.districts",
										"initialValue": []string{},
										"in": bson.M{
											"$concatArrays": []interface{}{
												"$$value",
												bson.M{"$cond": []interface{}{
													bson.M{"$regexMatch": bson.M{"input": "$$this.name", "regex": query, "options": "i"}},
													[]string{"$$this.name"},
													[]string{},
												}},
											},
										},
									}},
								},
							},
						}},
					},
				},
			},
		},
		{
			"$unwind": "$result",
		},
		{
			"$group": bson.M{
				"_id":     nil,
				"results": bson.M{"$addToSet": "$result"},
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
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []string{}, nil
	}

	return results[0].Results, nil
}
