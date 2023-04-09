package models

type User struct {
	UserName string `json:"userName,nonempty"`
	Email    string `json:"email,nonempty"`
	Password string `json:"password,nonempty"`
}

type Authentication struct {
	Email    string `json:"email,nonempty"`
	Password string `json:"password,nonempty"`
}
