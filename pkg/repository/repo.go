package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	conn *mongo.Client
}

func NewRepository(conn *mongo.Client) *Repository {
	return &Repository{
		conn: conn,
	}
}
