package models

type RoomInput struct {
	RoomId   int    `json:"room_id"`
	Password string `json:"password"`
	Private  bool   `json:"private"`
}

type Room struct {
	RoomId       int
	Private      bool
	Password     string
	CreatorEmail string
}

type JoinRoomInput struct {
	RoomId   int    `json:"room_id"`
	Password string `json:"password"`
}
