package repository

import (
	"github.com/gotalk/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	AddUser(user *models.User) error
	GetUserId(user *models.Authentication) (string, error)
}

type Repository struct {
	Authorization
}

func NewRepository(userCollection *mongo.Collection) *Repository {
	return &Repository{
		Authorization: NewUserRepository(userCollection),
	}
}
