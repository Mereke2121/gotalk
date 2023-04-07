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

	h.service.AddUser(&user)
}
