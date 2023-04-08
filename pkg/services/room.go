package services

import (
	"github.com/gotalk/models"
	"github.com/pkg/errors"
)

var rooms = make(map[int]models.Room)

func (s *Service) CreateRoom(input models.RoomInput, email string) (int, error) {
	if _, ok := rooms[input.RoomId]; !ok {
		room := models.Room{
			RoomId:       input.RoomId,
			Private:      input.Private,
			CreatorEmail: email,
		}
		if input.Private {
			room.Password = input.Password
		}
		rooms[input.RoomId] = room
		return input.RoomId, nil
	}
	return 0, errors.Errorf("room is already created; room id: %d", input.RoomId)
}
