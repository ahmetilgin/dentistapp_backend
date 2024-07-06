package mongo

import (
	"backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegionRepository struct {
    db             *mongo.Database
    regionCol     *mongo.Collection
    cityCol        *mongo.Collection
    districtCol    *mongo.Collection
}

func NewRegionRepository(db *mongo.Database, countriesCollection string, citiesCollection string, districtsCollection string) *RegionRepository {
    return &RegionRepository{
        db:             db,
        regionCol:     db.Collection(countriesCollection),
        cityCol:        db.Collection(citiesCollection),
        districtCol:    db.Collection(districtsCollection),
    }
}

func (r *RegionRepository) CreateRegion(ctx context.Context, region *models.Region) error {
	session, err := r.regionCol.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		for i := 0; i < len(region.Cities); i++ {
			city := &region.Cities[i]
			cityID := primitive.NewObjectID()
			city.ID = cityID

			for j := 0; j < len(city.Districts); j++ {
				district := &city.Districts[j]
				districtID := primitive.NewObjectID()
				district.ID = districtID
				district.CityID = cityID

				_, err := r.districtCol.InsertOne(sessCtx, district)
				if err != nil {
					return nil, err
				}
			}

			city.RegionID = region.ID
			_, err := r.cityCol.InsertOne(sessCtx, city)
			if err != nil {
				return nil, err
			}
		}

		regionID := primitive.NewObjectID()
		region.ID = regionID
		_, err := r.regionCol.InsertOne(sessCtx, region)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

func (r *RegionRepository) CreateCity(ctx context.Context, city *models.City) error {
    _, err := r.cityCol.InsertOne(ctx, city)
    return err
}

func (r *RegionRepository) CreateDistrict(ctx context.Context, district *models.District) error {
    _, err := r.districtCol.InsertOne(ctx, district)
    return err
}

func (r *RegionRepository) Search(ctx context.Context, query string) ([]string, error) {
    pipeline := []bson.M{
        {
            "$match": bson.M{
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
                            "input": "$cities",
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
                                        "input": "$$this.districts",
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
                "_id": nil,
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