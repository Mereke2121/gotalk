package handlers

import (
	"encoding/json"
	"github.com/gotalk/models"
	"go.uber.org/zap"
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
		h.logger.Error("parse input user model", zap.Error(err))
		handleError(http.StatusBadRequest, "parse input user model", w)
		return
	}

	err = h.service.AddUser(user)
	if err != nil {
		h.logger.Error("add user", zap.Error(err))
		handleError(http.StatusInternalServerError, "add user", w)
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
		h.logger.Error("parse input user model", zap.Error(err))
		handleError(http.StatusBadRequest, "parse input user model", w)
		return
	}

	token, err := h.service.Authenticate(&user)
	if err != nil {
		h.logger.Error("authenticate user", zap.Error(err))
		handleError(http.StatusInternalServerError, "authenticate user", w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
