package wallet

import "github.com/gin-gonic/gin"

func Register(rg *gin.RouterGroup) {
	h := NewHandler(NewService())
	wallets := rg.Group("/wallets")
	_ = h
	_ = wallets
}
