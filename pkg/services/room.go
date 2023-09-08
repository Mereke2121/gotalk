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

func NewRoomService(repo *repository.Repository) *RoomService {
	return &RoomService{
		rooms: make(map[int]*models.Room),
		repo:  repo,
	}
}

func (s *RoomService) CreateRoom(input *models.Room) (int, error) {
	// remove password if chat room is public
	if !input.Private {
		input.Password = ""
	}
	input.Password = hashPassword(input.Password)

	id, err := s.repo.Room.AddRoom(input)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *RoomService) UpdateRoom(input *models.UpdateRoomInput, roomId int, userId string) error {
	input.Password = hashPassword(input.Password)

	room, err := s.repo.GetRoomById(roomId)
	if err != nil {
		return err
	}
	if room.CreatorId != userId {
		return errors.New("unauthorized")
	}
	if !input.Private {
		input.Password = ""
	}

	err = s.repo.UpdateRoom(input, roomId)
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomService) AuthenticateInRoom(input *models.JoinRoomInput, roomId int, userId string) error {
	input.Password = hashPassword(input.Password)

	room, err := s.repo.GetRoomById(roomId)
	if err != nil {
		return err
	}
	if !room.Private && input == nil {
		return nil
	}
	if userId == room.CreatorId || room.Password == input.Password {
		return nil
	}

	return errors.Errorf("unauthorized for room id: %d", room.RoomId)
}

func (s *RoomService) GetAllRooms() ([]*models.RoomResponse, error) {
	rooms, err := s.repo.GetAllRooms()
	if err != nil {
		return nil, err
	}

	var result []*models.RoomResponse
	for _, room := range rooms {
		result = append(result, &models.RoomResponse{
			RoomId:    room.RoomId,
			Private:   room.Private,
			CreatorId: room.CreatorId,
		})
	}

	return result, nil
}

func (s *RoomService) GetRoomById(roomId int) (*models.RoomResponse, error) {
	room, err := s.repo.GetRoomById(roomId)
	if err != nil {
		return nil, err
	}
	return &models.RoomResponse{
		RoomId:    room.RoomId,
		Private:   room.Private,
		CreatorId: room.CreatorId,
	}, nil
}

func (s *RoomService) DeleteRoomById(roomId int, userId string) error {
	room, err := s.repo.GetRoomById(roomId)
	if err != nil {
		return err
	}
	if room.CreatorId != userId {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteRoomById(roomId)
}
