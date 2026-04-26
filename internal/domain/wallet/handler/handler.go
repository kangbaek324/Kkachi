package handler

import "github.com/kangbaek324/kkachi/internal/domain/wallet/service"

type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}
