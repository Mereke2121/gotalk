package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gotalk/models"
	"github.com/gotalk/utils"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) createRoom(w http.ResponseWriter, r *http.Request) {
	var room *models.Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		log.Println(err)
		return
	}

	// get user email from jwt token in header
	email, err := verifyUserEmail(r)

	roomId, err := h.service.CreateRoom(room, email)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(roomId)))
}

func (h *Handler) getAllRooms(w http.ResponseWriter, r *http.Request) {
	// verify user
	_, err := verifyUserEmail(r)
	if err != nil {
		log.Println(err)
		return
	}

	rooms, err := h.service.GetAllRooms()
	if err != nil {
		log.Println(err)
		return
	}

	resultBody, err := json.Marshal(rooms)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resultBody)
}

func (h *Handler) getRoomById(w http.ResponseWriter, r *http.Request) {
	// verify user
	_, err := verifyUserEmail(r)
	if err != nil {
		log.Println(err)
		return
	}

	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println("invalid room id")
		return
	}

	room, err := h.service.GetRoomById(roomId)
	if err != nil {
		log.Println(err)
		return
	}

	roomBody, err := json.Marshal(&room)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(roomBody)
}

func verifyUserEmail(r *http.Request) (string, error) {
	userToken, err := getJWTToken(r)
	if err != nil {
		return "", err
	}
	emailParam, err := utils.VerifyToken(userToken, utils.UserEmail)
	if err != nil {
		return "", err
	}
	email, ok := emailParam.(string)
	if !ok {
		log.Println("invalid user email")
	}
	return email, nil
}
