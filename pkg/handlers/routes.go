package handlers

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (h *Handler) InitRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/sign-up", h.signUp)
	mux.Post("/sign-in", h.signIn)

	mux.Get("/ws/{id}", h.wsConnection)
	mux.Post("/ws/{id}/join", h.joinRoom)

	return mux
}
