package services

import (
	"github.com/gotalk/models"
	"github.com/pkg/errors"
)

var rooms = make(map[int]models.Room)

func (s *Service) CreateRoom(input models.RoomInput, email string) error {
	if _, ok := rooms[input.RoomId]; !ok {
		rooms[input.RoomId] = models.Room{
			RoomId:       input.RoomId,
			Private:      input.Private,
			CreatorEmail: email,
		}
		return nil
	}
	return errors.Errorf("room is already created; room id: %d", input.RoomId)
}
