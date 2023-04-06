package handlers

import (
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var clients = make(map[int]map[*websocket.Conn]bool)

func (h *Handler) wsConnection(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	roomId := chi.URLParam(r, "id")
	if roomId == "" {
		log.Println("there is no provided room id")
		return
	}
	id, _ := strconv.Atoi(roomId)

	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if clients[id] == nil {
		clients[id] = make(map[*websocket.Conn]bool)
	}
	clients[id][conn] = true

	conn.WriteMessage(websocket.TextMessage, []byte("websocket connected"))

	go func(roomId int) {
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			// broadcast to all clients
			for client := range clients[roomId] {
				client.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}(id)
}
