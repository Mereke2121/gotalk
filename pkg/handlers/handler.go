package handlers

import "github.com/gotalk/pkg/service"

type Handler struct {
	service *service.Service
}

func NewHandler() *Handler {
	return &Handler{}
}
