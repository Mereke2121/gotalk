package services

import (
	"github.com/gorilla/websocket"
	"github.com/gotalk/models"
	"github.com/gotalk/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	AddUser(user *models.User) error
	Authenticate(user *models.Authentication) (string, error)
	GetUserById(userId string) (*models.User, error)
}

type Room interface {
	CreateRoom(input *models.Room) (int, error)
	UpdateRoom(input *models.UpdateRoomInput, roomId int, userId string) error
	AuthenticateInRoom(input *models.JoinRoomInput, roomId int, userId string) error
	GetAllRooms() ([]*models.RoomResponse, error)
	GetRoomById(roomId int) (*models.RoomResponse, error)
	DeleteRoomById(roomId int, userId string) error
}

type Websocket interface {
	MakeWSConnection(conn *websocket.Conn, roomId int, email string)
}

type Service struct {
	Authorization
	Room
	Websocket
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
		Room:          NewRoomService(repo),
		Websocket:     NewWebsocketService(),
	}
}
