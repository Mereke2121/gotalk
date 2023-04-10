package models

type User struct {
	UserName string `json:"username,nonempty"`
	UserId   string `json:"user_id,nonempty"`
	Password string `json:"password,nonempty"`
}

type Authentication struct {
	UserId   string `json:"user_id,nonempty"`
	Password string `json:"password,nonempty"`
}
