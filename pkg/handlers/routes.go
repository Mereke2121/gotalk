package handlers

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (h *Handler) InitRoutes() http.Handler {
	mux := chi.NewRouter()

	//auth
	mux.Post("/sign-up", h.signUp)
	mux.Post("/sign-in", h.signIn) // return jwt token which consists user email

	// rooms
	mux.Get("/room", h.getAllRooms)
	mux.Get("/room/{id}", h.getRoomById)
	mux.Put("/room/{id}", h.updateRoomById)
	mux.Delete("/room/{id}", h.deleteRoomById)
	mux.Post("/room", h.createRoom)

	// ws
	mux.Get("/ws/{id}", h.wsConnection)   // user id, room id from jwt token
	mux.Post("/ws/{id}/join", h.joinRoom) // for authorization inputs jwt token which consists email and room id

	return mux
}
