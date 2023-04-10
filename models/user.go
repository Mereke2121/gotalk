package models

type User struct {
	UserName string `json:"username,nonempty"`
	userId   string `json:"userId,nonempty"`
	Password string `json:"password,nonempty"`
}

type Authentication struct {
	userId   string `json:"userId,nonempty"`
	Password string `json:"password,nonempty"`
}
