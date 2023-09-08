package main

import (
	"context"
	"log"

	"github.com/gotalk/pkg/handlers"
	"github.com/gotalk/pkg/repository"
	"github.com/gotalk/pkg/services"
	"github.com/gotalk/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// @title 			Gotalk API
// @version			1.0
// @description		This is the chat rest api

// @host 			localhost:8080
// @BasePath		/
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// connect to mongo db
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.Fatal("connect to mongo db", zap.Error(err))
		return
	}

	// check mongo db connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.Fatal("ping mongo db connection", zap.Error(err))
		return
	}

	userCollection := client.Database("gotalk").Collection("users")
	roomCollection := client.Database("gotalk").Collection("rooms")

	repo := repository.NewRepository(userCollection, roomCollection)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service, logger)

	serverWS := new(server.WSServer)
	if err = serverWS.Run(":8000", handler.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}
