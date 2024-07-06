package mongo

import (
	"backend/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type JobRepository struct {
	jobCollection *mongo.Collection
	professionCollection *mongo.Collection
}

func NewJobRepository(jobCollection *mongo.Database, professionCollectionName, jobCollectionName string) *JobRepository {
	return &JobRepository{
		jobCollection: jobCollection.Collection(jobCollectionName),
		professionCollection: jobCollection.Collection(professionCollectionName),
	}
}

func (r JobRepository) CreateJob(ctx context.Context, user *models.BusinessUser, bm *models.Job) error {

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

func (r JobRepository) Search(ctx context.Context, location, keyword string) ([]*models.Job, error) {
	filter := bson.M{
        "$or": []bson.M{
            {"location": bson.M{"$regex": location, "$options": "i"}},
            {"job_title": bson.M{"$regex": keyword, "$options": "i"}},
            {"description": bson.M{"$regex": keyword, "$options": "i"}},
            {"requirements": bson.M{"$regex": keyword, "$options": "i"}},
        },
    }
	
    opts := options.Find().SetSort(bson.D{{Key: "date_posted", Value: -1}})

    cursor, err := r.jobCollection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var jobs []*models.Job
    for cursor.Next(ctx) {
        var job models.Job
        if err := cursor.Decode(&job); err != nil {
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


func (r JobRepository) SearchProfession(ctx context.Context, keyword string) ([]*models.Profession, error) {
    var results []*models.Profession

    // Filtre oluşturma
    filter := bson.M{"name": bson.M{"$regex": keyword, "$options": "i"}} // "i" opsiyonu, aramanın büyük/küçük harf duyarsız olmasını sağlar

    // Veritabanında arama yapma
    cursor, err := r.professionCollection.Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("error finding professions: %w", err)
    }
    defer cursor.Close(ctx) // Cursor'ı kapatmayı unutmayın

    // Sonuçları results dilimine yükleme
    if err = cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("error decoding professions: %w", err)
    }

    return results, nil
}


