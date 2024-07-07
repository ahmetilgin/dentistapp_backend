package mongo

import (
	"backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

type UserRepository struct {
	normalUserCollection *mongo.Collection
	businessUserCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, userCollectionString string, businessCollectionString string ) *UserRepository {
	return &UserRepository{
		normalUserCollection: db.Collection(userCollectionString),
		businessUserCollection: db.Collection(businessCollectionString),
	}
}

func (r UserRepository) CreateNormalUser(ctx context.Context, user *models.NormalUser) error {
	_, err := r.normalUserCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}


func (r UserRepository) CreateBusinessUser(ctx context.Context, user *models.BusinessUser) error {
	_, err := r.businessUserCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepository) GetNormalUser(ctx context.Context, username, password string)  (*models.NormalUser, error) {
	baseUser := new(models.NormalUser)
	err := r.normalUserCollection.FindOne(ctx, bson.M{
		"username": username,
		"password": password,
	}).Decode(baseUser)

	if err != nil {
		return nil, err
	}

	return baseUser, nil
}

func (r UserRepository) GetBusinessUser(ctx context.Context, username, password string)  (*models.BusinessUser, error) {
	baseUser := new(models.BusinessUser)
	err := r.businessUserCollection.FindOne(ctx, bson.M{
		"username": username,
		"password": password,
	}).Decode(baseUser)

	if err != nil {
		return nil, err
	}

	return baseUser, nil
}




