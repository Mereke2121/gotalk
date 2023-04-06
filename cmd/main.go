package main

import (
	"github.com/gotalk/pkg/handlers"
	"github.com/gotalk/server"
	"log"
)

func main() {
	handler := handlers.NewHandler()

	serverWS := new(server.WSServer)
	if err := serverWS.Run(":8000", handler.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}
