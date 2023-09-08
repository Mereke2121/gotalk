package repository

import (
	"context"

	"github.com/gotalk/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	usersCollection *mongo.Collection
}

func NewUserRepository(usersCollection *mongo.Collection) Authorization {
	return &UserRepository{usersCollection: usersCollection}
}

func (r *UserRepository) AddUser(user *models.User) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"username", -1},
			{"email", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.usersCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	_, err = r.usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.Wrap(err, "you are already authorized")
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetUserId(user *models.Authentication) (string, error) {
	var mongoUser *models.MongoUser
	err := r.usersCollection.FindOne(context.Background(), user).Decode(&mongoUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("you haven't authorized")
		}
		return "", err
	}

	return mongoUser.Id.Hex(), nil
}

func (r *UserRepository) GetUserById(userId string) (*models.User, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{"_id", id},
	}

	var user *models.User
	err = r.usersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
