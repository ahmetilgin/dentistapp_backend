package mongo

import (
	"backend/models"
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// escapeRegexSpecialChars escapes special characters in a string for use in a regular expression
func escapeRegexSpecialChars(input string) string {
	specialChars := []string{".", "^", "$", "*", "+", "?", "(", ")", "[", "]", "{", "}", "|", "\\"}
	for _, char := range specialChars {
		input = strings.ReplaceAll(input, char, "\\"+char)
	}
	return input
}

type JobRepository struct {
	jobCollection        *mongo.Collection
	professionCollection *mongo.Collection
}

func NewJobRepository(jobCollection *mongo.Database, professionCollectionName, jobCollectionName string) *JobRepository {
	return &JobRepository{
		jobCollection:        jobCollection.Collection(jobCollectionName),
		professionCollection: jobCollection.Collection(professionCollectionName),
	}
}

func (r JobRepository) CreateJob(ctx context.Context, user *models.BusinessUser, bm *models.Job) error {
	bm.UserID = user.ID
	_, err := r.jobCollection.InsertOne(ctx, bm)
	if err != nil {
		return err
	}

	return nil
}

func (r JobRepository) GetJobs(ctx context.Context) ([]*models.Job, error) {
	cur, err := r.jobCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	out := make([]*models.Job, 0)

	for cur.Next(ctx) {
		user := new(models.Job)
		err := cur.Decode(user)
		if err != nil {
			return nil, err
		}

		out = append(out, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r JobRepository) IncreaseSearchCounter(ctx context.Context, keyword, code string) (bool, error) {
	filterProfessions := bson.M{
		"name": keyword,
		"code": strings.ToUpper(code),
	}
	update := bson.M{"$inc": bson.M{"search_counter": 1}}
	errProfessions := r.professionCollection.FindOneAndUpdate(ctx, filterProfessions, update)
	if errProfessions == nil {
		err := errProfessions.Err()
		if err != nil {
			fmt.Printf("errProfessions.Err().Error(): %v\n", err.Error())
		}
		return false, err
	}

	return true, nil
}

func (r JobRepository) Search(ctx context.Context, location, keyword, region string) ([]*models.Job, error) {
	filter := bson.M{}

	if location != "-" {
		filter["location"] = bson.M{"$regex": location, "$options": "i"}
	}

	if keyword != "-" {
		filter["$or"] = []bson.M{
			{"job_title": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
			{"requirements": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "date_posted", Value: -1}})
	cursor, err := r.jobCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	ret, err := r.IncreaseSearchCounter(ctx, keyword, region)
	if !ret {
		if err != nil {
			fmt.Printf("err: %v\n", err.Error())
		}
		return nil, err
	}

	var jobs []*models.Job
	for cursor.Next(ctx) {
		var job models.Job
		if err := cursor.Decode(&job); err != nil {
			// print error
			fmt.Println(err.Error())
			return nil, err
		}
		jobs = append(jobs, &job)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r JobRepository) DeleteJob(ctx context.Context, user *models.BusinessUser, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := r.jobCollection.DeleteOne(ctx, bson.M{"_id": objID, "userId": user.ID})
	return err
}

func (r JobRepository) SearchProfession(ctx context.Context, keyword, code string) ([]*models.Profession, error) {
	var results []*models.Profession

	escapedKeyword := escapeRegexSpecialChars(keyword)

	filter := []bson.M{
		{"$match": bson.M{"code": strings.ToUpper(code)}},                                          // Code eşleşmesi
		{"$unwind": "$professions"},                                                                // Professions dizisini aç
		{"$match": bson.M{"professions.name": bson.M{"$regex": escapedKeyword, "$options": "i"}}},  // İsim filtreleme
		{"$project": bson.M{"name": "$professions.name", "count": "$professions.count", "_id": 0}}, // Sonuç olarak sadece name ve count al
		{"$sort": bson.M{"count": -1}},                                                             // count alanına göre azalan sırala
		{"$limit": 10},                                                                             // İlk 10 sonucu al
	}

	cursor, err := r.professionCollection.Aggregate(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding professions: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("error decoding professions: %w", err)
	}

	return results, nil
}

func (r JobRepository) GetPopulerJobs(ctx context.Context, code string) ([]*models.Profession, error) {
	filter := bson.M{
		"code": strings.ToUpper(code),
	}

	// Professions içindeki en popüler 5 mesleği almak için pipeline oluştur
	pipeline := []bson.M{
		{"$match": filter},
		{"$unwind": "$professions"},
		{"$sort": bson.M{"professions.count": -1}},
		{"$limit": 5},
		{"$project": bson.M{
			"name":  "$professions.name",
			"count": "$professions.count",
		}},
	}

	cursor, err := r.professionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var professions []*models.Profession
	for cursor.Next(ctx) {
		var profession models.Profession
		if err := cursor.Decode(&profession); err != nil {
			return nil, err
		}
		professions = append(professions, &profession)
	}

	// Cursor hatasını kontrol et
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return professions, nil
}
