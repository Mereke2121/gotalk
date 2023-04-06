package handlers

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (h *Handler) InitRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/ws/{id}", h.wsConnection)

	return mux
}
