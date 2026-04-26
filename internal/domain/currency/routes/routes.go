package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kangbaek324/kkachi/internal/domain/currency/handler"
	"github.com/kangbaek324/kkachi/internal/domain/currency/service"
)

func Register(rg *gin.RouterGroup) {
	h := handler.New(service.New())
	currencies := rg.Group("/currencies")
	_ = h
	_ = currencies
}
