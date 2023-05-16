package repository

import (
	"github.com/gotalk/api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	AddUser(user *models.User) error
	GetUserId(user *models.Authentication) (string, error)
	GetUserById(userId string) (*models.User, error)
}

type Room interface {
	AddRoom(input *models.Room) (int, error)
	GetAllRooms() ([]*models.Room, error)
	GetRoomById(roomId int) (*models.Room, error)
	UpdateRoom(input *models.UpdateRoomInput, roomId int) error
	DeleteRoomById(roomId int) error
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
