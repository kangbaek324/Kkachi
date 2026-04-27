package user

import "github.com/gin-gonic/gin"

func Register(rg *gin.RouterGroup) {
	h := NewHandler(NewService())
	users := rg.Group("/users")
	_ = h
	_ = users
}
