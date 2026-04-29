package wallet

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/kangbaek324/kkachi/db/sqlc"
	"github.com/kangbaek324/kkachi/internal/middleware"
)

func Register(rg *gin.RouterGroup, pool *pgxpool.Pool, jwtSecret string) {
	svc := NewService(db.New(pool))
	h := NewHandler(svc)

	wallets := rg.Group("/wallets")
	wallets.Use(middleware.Auth(jwtSecret))

	wallets.POST("/", h.createWallet)
}
