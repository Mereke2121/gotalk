package repository

import (
	"context"

	"github.com/gotalk/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoomRepository struct {
	roomCollection *mongo.Collection
}

func NewRoomRepository(roomCollection *mongo.Collection) *RoomRepository {
	return &RoomRepository{
		roomCollection: roomCollection,
	}
}

func (r *RoomRepository) AddRoom(input *models.Room) (int, error) {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"roomid", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.roomCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return 0, err
	}

	_, err = r.roomCollection.InsertOne(context.Background(), input)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return 0, errors.Errorf("you've already added room with this id: err: %s; id: %d", err.Error(), input.RoomId)
		}
		return 0, err
	}

	return input.RoomId, nil
}

func (r *RoomRepository) GetAllRooms() ([]*models.Room, error) {
	option := options.Find()

	cursor, err := r.roomCollection.Find(context.Background(), bson.D{}, option)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var rooms []*models.Room
	for cursor.Next(context.Background()) {
		var room *models.Room
		err := cursor.Decode(&room)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *RoomRepository) GetRoomById(roomId int) (*models.Room, error) {
	option := options.FindOne()
	filter := bson.D{
		{"roomid", roomId},
	}

	var room *models.Room
	err := r.roomCollection.FindOne(context.Background(), filter, option).Decode(&room)
	return room, err
}

func (r *RoomRepository) UpdateRoom(input *models.UpdateRoomInput, roomId int) error {
	filter := bson.D{
		{"roomid", roomId},
	}
	update := bson.D{
		{"$set", input},
	}

	return r.roomCollection.FindOneAndUpdate(context.Background(), filter, update).Err()
}

func (r *RoomRepository) DeleteRoomById(roomId int) error {
	filter := bson.D{
		{"roomid", roomId},
	}
	_, err := r.roomCollection.DeleteOne(context.Background(), filter)
	return err
}
