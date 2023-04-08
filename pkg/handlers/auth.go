package handlers

import (
	"encoding/json"
	"github.com/gotalk/models"
	"log"
	"net/http"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}
	if user.UserName == "" || user.Email == "" || user.Password == "" {
		log.Println("invalid body request")
		return
	}

	err = h.service.AddUser(&user)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authorized successfully"))
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var user models.Authentication
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}

	token, err := h.service.Authenticate(&user)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
