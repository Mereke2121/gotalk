package repository

import (
	"github.com/gotalk/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	AddUser(user *models.User) error
	GetUserId(user *models.Authentication) (string, error)
}

type Room interface {
	AddRoom(input *models.Room) (int, error)
	GetAllRooms() ([]*models.RoomResponse, error)
}

type Repository struct {
	Authorization
	Room
}

func NewRepository(userCollection, roomCollection *mongo.Collection) *Repository {
	return &Repository{
		Authorization: NewUserRepository(userCollection),
		Room:          NewRoomRepository(roomCollection),
	}
}
