package models

type Room struct {
	RoomId    int    `json:"room_id,nonempty"`
	Password  string `json:"password,omitempty"`
	Private   bool   `json:"private,nonempty"`
	CreatorId string `json:"creator_id,omitempty"`
}

type RoomResponse struct {
	RoomId    int    `json:"room_id"`
	Private   bool   `json:"private"`
	CreatorId string `json:"creator_id"`
}

type JoinRoomInput struct {
	Password string `json:"password,omitempty"`
}

type UpdateRoomInput struct {
	Password string `json:"password,omitempty"`
	Private  bool   `json:"private,nonempty"`
}
