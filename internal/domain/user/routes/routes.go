package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kangbaek324/kkachi/internal/domain/user/handler"
	"github.com/kangbaek324/kkachi/internal/domain/user/service"
)

func Register(rg *gin.RouterGroup) {
	h := handler.New(service.New())
	users := rg.Group("/users")
	_ = h
	_ = users
}
