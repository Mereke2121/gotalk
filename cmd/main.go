package main

import (
	"github.com/gotalk/pkg/handlers"
	"github.com/gotalk/pkg/services"
	"github.com/gotalk/server"
	"log"
)

func main() {
	service := services.NewService()
	handler := handlers.NewHandler(service)

	serverWS := new(server.WSServer)
	if err := serverWS.Run(":8000", handler.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}
