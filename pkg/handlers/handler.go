package handlers

import "github.com/gotalk/pkg/services"

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}
