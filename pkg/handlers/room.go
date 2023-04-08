package handlers

import (
	"encoding/json"
	"github.com/gotalk/models"
	"github.com/gotalk/utils"
	"log"
	"net/http"
)

func (h *Handler) createRoom(w http.ResponseWriter, r *http.Request) {
	var room models.RoomInput
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		log.Println(err)
		return
	}

	// get user email from jwt token in header
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

	err = h.service.CreateRoom(room, email)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
