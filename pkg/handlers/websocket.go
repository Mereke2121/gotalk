package handlers

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/gotalk/utils"
	"log"
	"net/http"
	"strconv"
)

var clients = make(map[int]map[*websocket.Conn]bool)

func (h *Handler) wsConnection(w http.ResponseWriter, r *http.Request) {
	// get body and header params
	roomId := chi.URLParam(r, "id")
	if roomId == "" {
		log.Println("there is no provided room id")
		return
	}
	id, _ := strconv.Atoi(roomId)

	// get jwt token from header
	token, err := getJWTToken(r)
	if err != nil {
		log.Println(err)
		return
	}

	// make authentication
	err = authenticate(id, token)
	if err != nil {
		log.Println(err)
		return
	}

	// make websocket connection
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

	// listen messages from websocket connection
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

func (h *Handler) joinRoom(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "id")
	if roomId == "" {
		log.Println("there is no provided room id")
		return
	}
	id, _ := strconv.Atoi(roomId)

	token, err := utils.CreateToken(id)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func getJWTToken(r *http.Request) (string, error) {
	token := r.Header.Get("token")
	if token == "" {
		return "", errors.New("there is no provided token")
	}
	return token, nil
}

func authenticate(roomId int, token string) error {
	id, err := utils.VerifyToken(token)
	if err != nil {
		return err
	}

	if roomId != id {
		return errors.New("you are unauthorized")
	}

	return nil
}
