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

type RoomResponse struct {
	RoomId       int    `json:"room_id"`
	Private      bool   `json:"private"`
	CreatorEmail string `json:"creator_email"`
}

type JoinRoomInput struct {
	RoomId   int    `json:"room_id"`
	Password string `json:"password"`
}
