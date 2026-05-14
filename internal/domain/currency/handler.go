package currency

import (
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

func (h *Handler) getCurrencies(c *gin.Context) {
	res, err := h.svc.getCurrencies(c.Request.Context())
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.ApiResponse(c, http.StatusOK, true, res)
}
