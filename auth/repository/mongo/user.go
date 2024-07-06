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
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collection string) *UserRepository {
	return &UserRepository{
		db: db.Collection(collection),
	}
}

func (r UserRepository) CreateNormalUser(ctx context.Context, user *models.NormalUser) error {
	_, err := r.db.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}


func (r UserRepository) CreateBusinessUser(ctx context.Context, user *models.BusinessUser) error {
	_, err := r.db.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepository) GetUser(ctx context.Context, username, password string)  (interface{}, error) {
	baseUser := new(models.BaseUser)
	err := r.db.FindOne(ctx, bson.M{
		"baseuser.username": username,
		"baseuser.password": password,
	}).Decode(baseUser)

	if err != nil {
		return nil, err
	}

	// Şimdi kullanıcı tipini belirle ve uygun struct'ı döndür
	var result bson.M
	err = r.db.FindOne(ctx, bson.M{"_id": baseUser.ID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	if _, hasBusiness := result["business_name"]; hasBusiness {
		var businessUser models.BusinessUser
		data, err := bson.Marshal(result)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(data, &businessUser)
		if err != nil {
			return nil, err
		}
		return businessUser, nil
	} else {
		var normalUser models.NormalUser
		data, err := bson.Marshal(result)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(data, &normalUser)
		if err != nil {
			return nil, err
		}
		return normalUser, nil
	}
}


