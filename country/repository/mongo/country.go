package mongo

import (
	"backend/models"
	"context"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func createCollation(code string) *options.Collation {
	var collation *options.Collation

	switch code {
	case "tr":
		collation = &options.Collation{
			Locale:   "tr",
			Strength: 1,
		}
	case "al":
		collation = &options.Collation{
			Locale:   "sq",
			Strength: 1,
		}
	case "en":
		collation = &options.Collation{
			Locale:   "en",
			Strength: 1,
		}
	default:
		collation = &options.Collation{
			Locale:   "en",
			Strength: 1,
		}
	}

	return collation
}

func (r *RegionRepository) Search(ctx context.Context, query, code string) ([]string, error) {
	if code == "en" {
		code = "al"
	}

	sanitizedQuery := regexp.QuoteMeta(query)

	pipeline := []bson.M{
		{"$match": bson.M{"code": strings.ToUpper(code)}},
		{"$unwind": "$cities"},
		{
			"$facet": bson.M{
				"cities_located": []bson.M{
					{"$match": bson.M{"cities.name": bson.M{"$regex": "^" + sanitizedQuery, "$options": "i"}}},
					{"$project": bson.M{"name": "$cities.name", "_id": 0}},
					{"$limit": 10},
				},
				"districts_located": []bson.M{
					{"$unwind": "$cities.districts"},
					{"$match": bson.M{"cities.districts.name": bson.M{"$regex": "^" + sanitizedQuery, "$options": "i"}}},
					{"$project": bson.M{"name": "$cities.districts.name", "_id": 0}},
					{"$limit": 10},
				},
			},
		},
		{
			"$project": bson.M{
				"results": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$gt": []interface{}{
								bson.M{"$size": "$cities_located"},
								0,
							},
						},
						"then": "$cities_located",
						"else": "$districts_located",
					},
				},
			},
		},
		{"$unwind": "$results"},
		{"$project": bson.M{"name": "$results.name"}},
		{"$limit": 10},
	}

	var results []struct {
		Name string `bson:"name"`
	}

	cursor, err := r.regionCol.Aggregate(ctx, pipeline, options.Aggregate().SetCollation(createCollation(code)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var cityAndDistricts []string
	for _, result := range results {
		cityAndDistricts = append(cityAndDistricts, result.Name)
	}

	return cityAndDistricts, nil
}
