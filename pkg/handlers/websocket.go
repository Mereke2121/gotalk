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
	userId, err := authenticate(roomId, token)
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

	go ListenWS(conn, roomId, userId)
}

func (h *Handler) joinRoom(w http.ResponseWriter, r *http.Request) {
	var input *models.JoinRoomInput
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Println(err)
			return
		}
	}

	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println("invalid room id")
		return
	}

	// get user id from jwt token in request header
	userId, err := verifyUserId(r)
	if err != nil {
		log.Println(err)
		return
	}

	// authentication for chat room by room id and user id
	err = h.service.AuthenticateInRoom(input, roomId, userId)
	if err != nil {
		log.Println(err)
		return
	}

	// create token for ws connection
	tokenWS, err := createToken(roomId, userId)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenWS))
}

func ListenWS(conn *websocket.Conn, roomId int, userId string) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		wsMessage := fmt.Sprintf("%s: %s", userId, string(msg))
		// broadcast to all clients
		for client := range clients[roomId] {
			client.WriteMessage(websocket.TextMessage, []byte(wsMessage))
		}
	}
}

func getJWTToken(r *http.Request) (string, error) {
	token := r.Header.Get("token")
	if token == "" {
		return "", errors.New("there is no provided token")
	}
	return token, nil
}

func createToken(roomId int, userId string) (string, error) {
	jwtFields := []utils.JWTTokenField{
		{
			Type:  utils.RoomId,
			Value: roomId,
		},
		{
			Type:  utils.UserId,
			Value: userId,
		},
	}
	tokenWS, err := utils.CreateToken(jwtFields)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenWS, nil
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

	userIdParam, err := utils.VerifyToken(token, utils.UserId)
	if err != nil {
		return "", err
	}
	userId, ok := userIdParam.(string)
	if !ok {
		return "", errors.New("convert user userId in jwt token from interface to string")
	}

	return userId, nil
}
