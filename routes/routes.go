package routes

import (
	"github.com/gin-gonic/gin"
	currencyroutes "github.com/kangbaek324/kkachi/internal/domain/currency/routes"
	userroutes "github.com/kangbaek324/kkachi/internal/domain/user/routes"
	walletroutes "github.com/kangbaek324/kkachi/internal/domain/wallet/routes"
)

func Register(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	userroutes.Register(v1)
	walletroutes.Register(v1)
	currencyroutes.Register(v1)
}
