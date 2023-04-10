package services

import (
	"github.com/gotalk/models"
	"github.com/gotalk/pkg/repository"
)

type Authorization interface {
	AddUser(user *models.User) error
	Authenticate(user *models.Authentication) (string, error)
}

type Room interface {
	CreateRoom(input *models.Room) (int, error)
	UpdateRoom(input *models.UpdateRoomInput, roomId int, userId string) error
	AuthenticateInRoom(input *models.JoinRoomInput, roomId int, userId string) error
	GetAllRooms() ([]models.RoomResponse, error)
	GetRoomById(roomId int) (*models.RoomResponse, error)
	DeleteRoomById(roomId int, userId string) error
}

type Service struct {
	Authorization
	Room
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
		Room:          NewRoomService(repo),
	}
}
