package wallet

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

func (h *Handler) createWallet(c *gin.Context) {
	req, ok := common.BindJSON[CreateWalletRequest](c)
	if !ok {
		return
	}

	res, err := h.svc.CreateWallet(c.Request.Context(), req, c.GetInt64("userId"))
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.ApiResponse(c, http.StatusCreated, true, res)
}

func (h *Handler) getWallets(c *gin.Context) {
	res, err := h.svc.GetWallets(c.Request.Context(), c.GetInt64("userId"))
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.ApiResponse(c, http.StatusOK, true, res)
}

func (h *Handler) editWalletNickname(c *gin.Context) {
	req, ok := common.BindJSON[EditWalletNicknameRequest](c)
	if !ok {
		return
	}

	if err := h.svc.EditWalletNickname(c.Request.Context(), req, c.GetInt64("userId")); err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.ApiResponse(c, http.StatusOK, true, nil)
}
