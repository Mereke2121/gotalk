package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/gotalk/models"
	"github.com/gotalk/utils"
	"log"
	"net/http"
	"strconv"
)

var clients = make(map[int]map[*websocket.Conn]bool)

func (h *Handler) wsConnection(w http.ResponseWriter, r *http.Request) {
	// get body and header params
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println("invalid room id")
		return
	}

	// get jwt token from header
	token, err := getJWTToken(r)
	if err != nil {
		log.Println(err)
		return
	}

	// make authentication
	email, err := authenticate(roomId, token)
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
	if clients[roomId] == nil {
		clients[roomId] = make(map[*websocket.Conn]bool)
	}
	clients[roomId][conn] = true

	// listen messages from websocket connection
	conn.WriteMessage(websocket.TextMessage, []byte("websocket connected"))

	go func(roomId int, userEmail string) {
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			wsMessage := fmt.Sprintf("%s: %s", userEmail, string(msg))
			// broadcast to all clients
			for client := range clients[roomId] {
				client.WriteMessage(websocket.TextMessage, []byte(wsMessage))
			}
		}
	}(roomId, email)
}

func (h *Handler) joinRoom(w http.ResponseWriter, r *http.Request) {
	var input *models.JoinRoomInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println(err)
		return
	}
	roomId := input.RoomId

	// get user email from jwt token in request header
	userToken, err := getJWTToken(r)
	if err != nil {
		log.Println(err)
		return
	}

	emailParam, err := utils.VerifyToken(userToken, utils.UserEmail)
	if err != nil {
		log.Println(err)
		return
	}
	email, ok := emailParam.(string)
	if !ok {
		log.Println("invalid user email")
	}

	// authentication for chat room by room id and user email
	err = h.service.AuthenticateInRoom(input, email)
	if err != nil {
		log.Println(err)
		return
	}

	// create token for ws connection
	jwtFields := []utils.JWTTokenField{
		{
			Type:  utils.RoomId,
			Value: roomId,
		},
		{
			Type:  utils.UserEmail,
			Value: email,
		},
	}
	tokenWS, err := utils.CreateToken(jwtFields)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenWS))
}

func getJWTToken(r *http.Request) (string, error) {
	token := r.Header.Get("token")
	if token == "" {
		return "", errors.New("there is no provided token")
	}
	return token, nil
}

func authenticate(roomIdHeader int, token string) (string, error) {
	roomParam, err := utils.VerifyToken(token, utils.RoomId)
	if err != nil {
		return "", err
	}
	roomId, ok := roomParam.(float64)
	if !ok {
		return "", errors.New("convert room id in jwt token from interface to int")
	}

	if roomIdHeader != int(roomId) {
		return "", errors.New("you are unauthorized")
	}

	emailParam, err := utils.VerifyToken(token, utils.UserEmail)
	if err != nil {
		return "", err
	}
	email, ok := emailParam.(string)
	if !ok {
		return "", errors.New("convert user email in jwt token from interface to string")
	}

	return email, nil
}
