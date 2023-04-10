package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoUser struct {
	Id       primitive.ObjectID `bson:"_id"`
	UserName string             `json:"username"`
	userId   string             `json:"userId"`
	Password string             `json:"password"`
}
