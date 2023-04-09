package models

type Room struct {
	RoomId       int    `json:"room_id,nonempty"`
	Password     string `json:"password,omitempty"`
	Private      bool   `json:"private,nonempty"`
	CreatorEmail string `json:"creator_email,omitempty"`
}

type RoomResponse struct {
	RoomId       int    `json:"room_id"`
	Private      bool   `json:"private"`
	CreatorEmail string `json:"creator_email"`
}

type JoinRoomInput struct {
	Password string `json:"password,omitempty"`
}

type UpdateRoomInput struct {
	Password string `json:"password,omitempty"`
	Private  bool   `json:"private,nonempty"`
}
