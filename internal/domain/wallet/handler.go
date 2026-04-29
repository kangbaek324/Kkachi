package wallet

import (
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

func (h *Handler) createWallet(c *gin.Context) {
	req, ok := common.BindJSON[CreateWalletRequest](c)
	if !ok {
		return
	}

	res, err := h.svc.CreateWallet(c.Request.Context(), req, c.GetInt64("userId"))
	if err != nil {
		log.Printf("createWallet: %v", err)
		common.ApiResponse(c, http.StatusInternalServerError, false, nil, "internal server error")
		return
	}
	common.ApiResponse(c, http.StatusCreated, true, res)
}

func (h *Handler) getWallets(c *gin.Context) {
	res, err := h.svc.GetWallets(c.Request.Context(), c.GetInt64("userId"))
	if err != nil {
		log.Printf("createWallet: %v", err)
		common.ApiResponse(c, http.StatusInternalServerError, false, nil, "internal server error")
		return
	}
	common.ApiResponse(c, http.StatusCreated, true, res)
}
