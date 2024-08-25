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
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type UserRepository struct {
	normalUserCollection         *mongo.Collection
	businessUserCollection       *mongo.Collection
	passwordResetTokenCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, userCollectionString string, businessCollectionString string, passwordResetTokenCollectioNString string) *UserRepository {
	return &UserRepository{
		normalUserCollection:         db.Collection(userCollectionString),
		businessUserCollection:       db.Collection(businessCollectionString),
		passwordResetTokenCollection: db.Collection(passwordResetTokenCollectioNString),
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

func (r UserRepository) GetNormalUser(ctx context.Context, email, password string) (*models.NormalUser, error) {
	baseUser := new(models.NormalUser)
	err := r.normalUserCollection.FindOne(ctx, bson.M{
		"email":    email,
		"password": password,
	}).Decode(baseUser)

	if err != nil {
		return nil, err
	}

	return baseUser, nil
}

func (r UserRepository) GetBusinessUser(ctx context.Context, email, password string) (*models.BusinessUser, error) {
	baseUser := new(models.BusinessUser)
	err := r.businessUserCollection.FindOne(ctx, bson.M{
		"email":    email,
		"password": password,
	}).Decode(baseUser)

	if err != nil {
		return nil, err
	}

	return baseUser, nil
}

func (r UserRepository) GetNormalUserByEmail(ctx context.Context, email string) (*models.NormalUser, error) {
	baseUser := new(models.NormalUser)
	err := r.normalUserCollection.FindOne(ctx, bson.M{"email": email}).Decode(baseUser)
	if err != nil {
		return nil, err
	}
	return baseUser, nil
}

func (r UserRepository) GetBusinessUserByEmail(ctx context.Context, email string) (*models.BusinessUser, error) {
	baseUser := new(models.BusinessUser)
	err := r.businessUserCollection.FindOne(ctx, bson.M{"email": email}).Decode(baseUser)
	if err != nil {
		return nil, err
	}
	return baseUser, nil
}

func (r UserRepository) InsetPasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error {
	_, err := r.passwordResetTokenCollection.InsertOne(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepository) CheckPasswordResetToken(ctx context.Context, userID primitive.ObjectID, token string) error {
	result := r.passwordResetTokenCollection.FindOne(ctx, bson.M{"user_id": userID, "token": token})
	err := result.Err()
	if err != nil {
		return err
	}

	_, err = r.passwordResetTokenCollection.DeleteOne(ctx, bson.M{"user_id": userID, "token": token})
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) UpdatePassword(ctx context.Context, user interface{}, token, newPassword string) error {
	var userID primitive.ObjectID

	if normalUser, ok := user.(*models.NormalUser); ok {
		userID = normalUser.ID
		err := r.CheckPasswordResetToken(ctx, userID, token)
		if err != nil {
			return err
		}

		_, err = r.normalUserCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"password": newPassword}})
		if err != nil {
			return err
		}
	} else if businessUser, ok := user.(*models.BusinessUser); ok {
		userID = businessUser.ID
		err := r.CheckPasswordResetToken(ctx, userID, token)
		if err != nil {
			return err
		}
		_, err = r.businessUserCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"password": newPassword}})
		if err != nil {
			return err
		}
	}
	return nil
}
