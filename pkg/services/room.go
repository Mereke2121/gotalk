package services

import (
	"github.com/gotalk/models"
	"github.com/pkg/errors"
)

var rooms = make(map[int]models.Room)

func (s *Service) CreateRoom(input *models.RoomInput, email string) (int, error) {
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

func (s *Service) AuthenticateInRoom(input *models.JoinRoomInput, email string) error {
	room, ok := rooms[input.RoomId]
	if !ok {
		return errors.Errorf("there is no room by id: %d", input.RoomId)
	}
	if !room.Private || email == room.CreatorEmail || room.Password == input.Password {
		return nil
	}
	return errors.Errorf("unauthorized for room id: %d", input.RoomId)
}

func (s *Service) GetAllRooms() ([]models.RoomResponse, error) {
	var result []models.RoomResponse

	for _, room := range rooms {
		result = append(result, models.RoomResponse{
			RoomId:       room.RoomId,
			Private:      room.Private,
			CreatorEmail: room.CreatorEmail,
		})
	}

	return result, nil
}
