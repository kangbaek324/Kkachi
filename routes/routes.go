package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kangbaek324/kkachi/internal/common"
	"github.com/kangbaek324/kkachi/internal/domain/currency"
	"github.com/kangbaek324/kkachi/internal/domain/user"
	"github.com/kangbaek324/kkachi/internal/domain/wallet"
)

func Register(r *gin.Engine, pool *pgxpool.Pool) {
	r.GET("/health", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, common.Response{
				Code:    http.StatusServiceUnavailable,
				Success: false,
				Message: "database unavailable",
			})
			return
		}
		c.JSON(http.StatusOK, common.Response{
			Code:    http.StatusOK,
			Success: true,
			Message: "ok",
		})
	})

	v1 := r.Group("/api/v1")
	user.Register(v1)
	wallet.Register(v1)
	currency.Register(v1)
}
