package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/kangbaek324/kkachi/db/sqlc"
)

func Register(rg *gin.RouterGroup, pool *pgxpool.Pool, jwtSecret string) {
	svc := NewService(db.New(pool), jwtSecret)
	h := NewHandler(svc)

	users := rg.Group("/users")

	users.POST("/register", h.register)
	users.POST("/login", h.login)
	users.POST("/refresh-accesstoken", h.refreshAccessToken)
}
