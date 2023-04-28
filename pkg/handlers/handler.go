package handlers

import (
	"github.com/gotalk/pkg/services"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	service *services.Service
	logger  *zap.Logger
}

func NewHandler(service *services.Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func handleError(code int, message string, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}
