package mongo

import (
	"backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



type JobRepository struct {
	db *mongo.Collection
}

func NewJobRepository(db *mongo.Database, collection string) *JobRepository {
	return &JobRepository{
		db: db.Collection(collection),
	}
}

func (r JobRepository) CreateJob(ctx context.Context, user *models.User, bm *models.Job) error {

	res, err := r.db.InsertOne(ctx, bm)
	if err != nil {
		return err
	}

	bm.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r JobRepository) GetJobs(ctx context.Context) ([]*models.Job, error) {
	cur, err := r.db.Find(ctx, bson.M{})

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

func (r JobRepository) DeleteJob(ctx context.Context, user *models.User, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	uID, _ := primitive.ObjectIDFromHex(user.ID)

	_, err := r.db.DeleteOne(ctx, bson.M{"_id": objID, "userId": uID})
	return err
}



