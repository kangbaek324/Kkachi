package user

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kangbaek324/kkachi/internal/common"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) register(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiResponse(c, http.StatusBadRequest, false, nil, err.Error())
		return
	}

	result, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUsernameAlreadyExists) {
			common.ApiResponse(c, http.StatusConflict, false, nil, "username already exists")
			return
		}
		log.Printf("register: %v", err)
		common.ApiResponse(c, http.StatusInternalServerError, false, nil, "internal server error")
		return
	}

	common.ApiResponse(c, http.StatusCreated, true, result)
}
