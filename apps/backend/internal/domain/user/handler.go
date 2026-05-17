package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kangbaek324/kkachi/apps/backend/internal/common"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) register(c *gin.Context) {
	req, ok := common.BindJSON[CreateUserRequest](c)
	if !ok {
		return
	}

	res, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.ApiResponse(c, http.StatusCreated, true, res)
}

func (h *Handler) login(c *gin.Context) {
	req, ok := common.BindJSON[LoginRequest](c)
	if !ok {
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.ApiResponse(c, http.StatusOK, true, res)
}

func (h *Handler) refreshAccessToken(c *gin.Context) {
	req, ok := common.BindJSON[RefreshAccessTokenRequest](c)
	if !ok {
		return
	}

	res, err := h.svc.RefreshAccessToken(c.Request.Context(), req)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.ApiResponse(c, http.StatusOK, true, res)
}
