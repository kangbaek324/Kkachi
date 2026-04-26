package handler

import "github.com/kangbaek324/kkachi/internal/domain/user/service"

type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}
