package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kangbaek324/kkachi/internal/domain/wallet/handler"
	"github.com/kangbaek324/kkachi/internal/domain/wallet/service"
)

func Register(rg *gin.RouterGroup) {
	h := handler.New(service.New())
	wallets := rg.Group("/wallets")
	_ = h
	_ = wallets
}
