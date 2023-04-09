package services

import (
	"github.com/gotalk/models"
	"github.com/pkg/errors"
)

var rooms = make(map[int]models.Room)

func (s *Service) CreateRoom(input *models.Room, email string) (int, error) {
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

func (s *Service) AuthenticateInRoom(input *models.JoinRoomInput, roomId int, email string) error {
	room, ok := rooms[roomId]
	if !ok {
		return errors.Errorf("there is no room by id: %d", roomId)
	}
	if !room.Private || email == room.CreatorEmail || room.Password == input.Password {
		return nil
	}
	return errors.Errorf("unauthorized for room id: %d", room)
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

func (s *Service) GetRoomById(roomId int) (*models.RoomResponse, error) {
	room, ok := rooms[roomId]
	if !ok {
		return nil, errors.Errorf("there's no room for provided room id: %d", roomId)
	}
	return &models.RoomResponse{
		RoomId:       room.RoomId,
		Private:      room.Private,
		CreatorEmail: room.CreatorEmail,
	}, nil
}
