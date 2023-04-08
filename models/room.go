package models

type RoomInput struct {
	Private bool `json:"is_private"`
	RoomId  int  `json:"room_id"`
}

type Room struct {
	RoomId       int
	Private      bool
	CreatorEmail string
}
