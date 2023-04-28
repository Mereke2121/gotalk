package handlers

import (
	"fmt"
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

func handleError(code int, message string, err error, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("%s; err: %s", message, err.Error())))
}
