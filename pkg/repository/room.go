package repository

import (
	"github.com/gotalk/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomRepository struct {
	roomCollection *mongo.Collection
}

func NewRoomRepository(roomCollection *mongo.Collection) *RoomRepository {
	return &RoomRepository{
		roomCollection: roomCollection,
	}
}

func (r *RoomRepository) AddRoom(input *models.Room, userId string) (int, error) {
	return 0, nil
}
