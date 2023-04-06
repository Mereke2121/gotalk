package handlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)

func (h *Handler) wsConnection(w http.ResponseWriter, r *http.Request) {
	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn.WriteJSON("websocket connected")
	clients[conn] = true

	go func() {
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			// broadcast to all clients
			for client := range clients {
				client.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}()
}
