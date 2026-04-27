package currency

import "github.com/gin-gonic/gin"

func Register(rg *gin.RouterGroup) {
	h := NewHandler(NewService())
	currencies := rg.Group("/currencies")
	_ = h
	_ = currencies
}
