package currency

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/kangbaek324/kkachi/apps/backend/db/sqlc"
)

func Register(rg *gin.RouterGroup, pool *pgxpool.Pool) {
	h := NewHandler(NewService(db.New(pool)))
	currencies := rg.Group("/currencies")

	currencies.GET("/", h.getCurrencies)
}
