package models

type Room struct {
	RoomId       int    `json:"room_id,nonempty"`
	Password     string `json:"password,omitempty"`
	Private      bool   `json:"private,nonempty"`
	CreatorEmail string `json:"creator_email,omitempty"`
}

type RoomResponse struct {
	RoomId       int
	Private      bool
	CreatorEmail string
}

type JoinRoomInput struct {
	Password string `json:"password,omitempty"`
}
