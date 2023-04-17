package handlers

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	_ "github.com/gotalk/docs"
)

func (h *Handler) InitRoutes() http.Handler {
	mux := chi.NewRouter()

	// swagger
	mux.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/swagger/doc.json")))

	// auth
	mux.Post("/sign-up", h.signUp)
	mux.Post("/sign-in", h.signIn) // return jwt token which consists user id

	// rooms
	mux.Get("/room", h.getAllRooms)
	mux.Post("/room", h.createRoom)
	mux.Get("/room/{id}", h.getRoomById)
	mux.Put("/room/{id}", h.updateRoomById)
	mux.Delete("/room/{id}", h.deleteRoomById)

	// ws
	mux.Get("/ws/{id}", h.wsConnection)   // user id, room id from jwt token
	mux.Post("/ws/{id}/join", h.joinRoom) // for authorization inputs jwt token which consists user id and room id

	return mux
}
