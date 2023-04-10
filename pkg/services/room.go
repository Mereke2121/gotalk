package services

import (
	"github.com/gotalk/models"
	"github.com/gotalk/pkg/repository"
	"github.com/pkg/errors"
)

type RoomService struct {
	repo  *repository.Repository
	rooms map[int]*models.Room
}

func NewRoomService(repo *repository.Repository) Room {
	return &RoomService{
		rooms: make(map[int]*models.Room),
		repo:  repo,
	}
}

func (s *RoomService) CreateRoom(input *models.Room, userId string) (int, error) {
	if _, ok := s.rooms[input.RoomId]; !ok {
		room := models.Room{
			RoomId:        input.RoomId,
			Private:       input.Private,
			CreatoruserId: userId,
		}
		if input.Private {
			room.Password = input.Password
		}
		s.rooms[input.RoomId] = &room
		return input.RoomId, nil
	}
	return 0, errors.Errorf("room is already created; room id: %d", input.RoomId)
}

func (s *RoomService) UpdateRoom(input *models.UpdateRoomInput, roomId int, userId string) error {
	room, ok := s.rooms[roomId]
	if !ok {
		return errors.Errorf("there is no room for by provided room id: %d", room.RoomId)
	}
	if room.CreatoruserId != userId {
		return errors.New("unauthorized")
	}

	room.Private = input.Private
	if input.Private {
		room.Password = input.Password
	}
	room.Password = ""

	return nil
}

func (s *RoomService) AuthenticateInRoom(input *models.JoinRoomInput, roomId int, userId string) error {
	room, ok := s.rooms[roomId]
	if !ok {
		return errors.Errorf("there is no room by id: %d", roomId)
	}
	if !room.Private && input == nil {
		return nil
	}
	if userId == room.CreatoruserId || room.Password == input.Password {
		return nil
	}
	return errors.Errorf("unauthorized for room id: %d", room)
}

func (s *RoomService) GetAllRooms() ([]models.RoomResponse, error) {
	var result []models.RoomResponse

	for _, room := range s.rooms {
		result = append(result, models.RoomResponse{
			RoomId:        room.RoomId,
			Private:       room.Private,
			CreatoruserId: room.CreatoruserId,
		})
	}

	return result, nil
}

func (s *RoomService) GetRoomById(roomId int) (*models.RoomResponse, error) {
	room, ok := s.rooms[roomId]
	if !ok {
		return nil, errors.Errorf("there's no room for provided room id: %d", roomId)
	}
	return &models.RoomResponse{
		RoomId:        room.RoomId,
		Private:       room.Private,
		CreatoruserId: room.CreatoruserId,
	}, nil
}

func (s *RoomService) DeleteRoomById(roomId int, userId string) error {
	room, ok := s.rooms[roomId]
	if !ok {
		return errors.Errorf("there is no room for provided room id: %d", roomId)
	}
	if room.CreatoruserId != userId {
		return errors.New("you don't have access for delete")
	}

	delete(s.rooms, roomId)
	return nil
}
