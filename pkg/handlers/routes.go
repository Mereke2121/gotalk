package handlers

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (h *Handler) InitRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/sign-up", h.signUp)
	mux.Post("/sign-in", h.signIn) // return jwt token which consists user email

	mux.Post("/ws/room", h.createRoom) // creates room
	mux.Get("/ws/room", h.getAllRooms)

	mux.Get("/ws/{id}", h.wsConnection)   // user id, room id from jwt token
	mux.Post("/ws/{id}/join", h.joinRoom) // for authorization inputs jwt token which consists email and room id

	return mux
}
