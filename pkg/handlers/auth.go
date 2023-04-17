package handlers

import (
	"encoding/json"
	"github.com/gotalk/models"
	"log"
	"net/http"
)

// @Summary      Sign up
// @Description  authorization
// @Tags         auth
// @Param 		 input body models.User true "Authorization"
// @Success      200  {string}  string
// @Router       /sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.service.AddUser(user)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authorized successfully"))
}

// @Summary      Sign in
// @Description  authentication
// @Tags         auth
// @Param        input body models.Authentication true "Authentication"
// @Success      200  {string} token
// @Router       /sign-in [post]
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
