package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gotalk/models"
	"github.com/gotalk/utils"
	"github.com/pkg/errors"
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
	if room.Private && room.Password == "" {
		log.Println(errors.New("there is no provided password"))
	}

	// get user userId from jwt token in header
	userId, err := verifyUserId(r)
	if err != nil {
		log.Println(err)
		return
	}
	room.CreatorId = userId

	roomId, err := h.service.CreateRoom(room)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(roomId)))
}

func (h *Handler) updateRoomById(w http.ResponseWriter, r *http.Request) {
	var room *models.UpdateRoomInput
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		log.Println(err)
		return
	}
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println("invalid room id")
		return
	}

	// get user id from jwt token in header
	userId, err := verifyUserId(r)

	err = h.service.UpdateRoom(room, roomId, userId)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully updated the room"))
}

func (h *Handler) getAllRooms(w http.ResponseWriter, r *http.Request) {
	// verify user
	_, err := verifyUserId(r)
	if err != nil {
		log.Println(err)
		return
	}

	rooms, err := h.service.GetAllRooms()
	if err != nil {
		log.Println(err)
		return
	}

	resultBody, err := json.MarshalIndent(rooms, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resultBody)
}

func (h *Handler) getRoomById(w http.ResponseWriter, r *http.Request) {
	// verify user
	_, err := verifyUserId(r)
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

func (h *Handler) deleteRoomById(w http.ResponseWriter, r *http.Request) {
	// verify user
	userId, err := verifyUserId(r)
	if err != nil {
		log.Println(err)
		return
	}

	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println("invalid room id")
		return
	}

	err = h.service.DeleteRoomById(roomId, userId)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete successfully"))
}

func verifyUserId(r *http.Request) (string, error) {
	userToken, err := getJWTToken(r)
	if err != nil {
		return "", err
	}
	userIdStr, err := utils.VerifyToken(userToken, utils.UserId)
	if err != nil {
		return "", err
	}
	userId, ok := userIdStr.(string)
	if !ok {
		log.Println("invalid user id")
	}
	return userId, nil
}
