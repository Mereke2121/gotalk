package models

type Room struct {
	RoomId        int    `json:"room_id,nonempty"`
	Password      string `json:"password,omitempty"`
	Private       bool   `json:"private,nonempty"`
	CreatoruserId string `json:"creator_userId,omitempty"`
}

type RoomResponse struct {
	RoomId        int    `json:"room_id"`
	Private       bool   `json:"private"`
	CreatoruserId string `json:"creator_userId"`
}

type JoinRoomInput struct {
	Password string `json:"password,omitempty"`
}

type UpdateRoomInput struct {
	Password string `json:"password,omitempty"`
	Private  bool   `json:"private,nonempty"`
}
