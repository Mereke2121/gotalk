package repository

import (
	"context"
	"github.com/gotalk/models"
	"github.com/gotalk/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Repository struct {
	conn            *mongo.Client
	usersCollection *mongo.Collection
}

func NewRepository(conn *mongo.Client) *Repository {
	userCollection := conn.Database("gotalk").Collection("users")
	return &Repository{
		usersCollection: userCollection,
		conn:            conn,
	}
}

func (r *Repository) AddUser(user *models.User) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{string(utils.UserName), -1},
			{string(utils.UserEmail), 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.usersCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	res, err := r.usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.Wrap(err, "you are already authorized")
		}
		return err
	}

	log.Println("inserted user with id: ", res.InsertedID)
	return nil
}
