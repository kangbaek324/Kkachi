package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/kangbaek324/kkachi/db/sqlc"
)

func Register(rg *gin.RouterGroup, pool *pgxpool.Pool) {
	svc := NewService(db.New(pool))
	h := NewHandler(svc)

	users := rg.Group("/users")
	users.POST("/register", h.register)
}
