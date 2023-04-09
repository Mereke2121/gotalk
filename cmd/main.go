package main

import (
	"context"
	"github.com/gotalk/pkg/handlers"
	"github.com/gotalk/pkg/repository"
	"github.com/gotalk/pkg/services"
	"github.com/gotalk/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// connect to mongo db
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return
	}

	// check mongo db connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	repo := repository.NewRepository(client)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)

	serverWS := new(server.WSServer)
	if err = serverWS.Run(":8000", handler.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}
