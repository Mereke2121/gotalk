package services

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	clients map[int]map[*websocket.Conn]bool
}

func NewWebsocketService() Websocket {
	clients := make(map[int]map[*websocket.Conn]bool)
	return &WebsocketService{clients: clients}
}

func (s *WebsocketService) MakeWSConnection(conn *websocket.Conn, roomId int, email string) {
	if s.clients[roomId] == nil {
		s.clients[roomId] = make(map[*websocket.Conn]bool)
	}
	s.clients[roomId][conn] = true

	// listen messages from websocket connection
	conn.WriteMessage(websocket.TextMessage, []byte("websocket connected"))

	go s.listenWS(conn, roomId, email)
}

func (s *WebsocketService) listenWS(conn *websocket.Conn, roomId int, email string) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		wsMessage := fmt.Sprintf("%s: %s", email, string(msg))
		// broadcast to all clients
		for client := range s.clients[roomId] {
			client.WriteMessage(websocket.TextMessage, []byte(wsMessage))
		}
	}
}
