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

	res, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUsernameAlreadyExists) {
			common.ApiResponse(c, http.StatusConflict, false, nil, err.Error())
			return
		}
		log.Printf("register: %v", err)
		common.ApiResponse(c, http.StatusInternalServerError, false, nil, "Internal server error")
		return
	}

	common.ApiResponse(c, http.StatusCreated, true, res)
}

func (h *Handler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiResponse(c, http.StatusBadRequest, false, nil, err.Error())
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			common.ApiResponse(c, http.StatusUnauthorized, false, nil, err.Error())
			return
		}
		log.Printf("login: %v", err)
		common.ApiResponse(c, http.StatusInternalServerError, false, nil, "Internal server error")
		return
	}

	common.ApiResponse(c, http.StatusOK, true, res)
}

func (h *Handler) refreshAccessToken(c *gin.Context) {
	var req RefreshAccessTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiResponse(c, http.StatusBadRequest, false, nil, err.Error())
		return
	}

	res, err := h.svc.RefreshAccessToken(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			common.ApiResponse(c, http.StatusUnauthorized, false, nil, err.Error())
			return
		}
		log.Printf("refreshAccessToken: %v", err)
		common.ApiResponse(c, http.StatusInternalServerError, false, nil, "Internal server error")
		return
	}

	common.ApiResponse(c, http.StatusOK, true, res)
}
