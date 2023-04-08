package models

type RoomInput struct {
	RoomId   int    `json:"room_id"`
	Password string `json:"password"`
	Private  bool   `json:"is_private"`
}

type Room struct {
	RoomId       int
	Private      bool
	Password     string
	CreatorEmail string
}
