package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kangbaek324/kkachi/internal/common"
	"github.com/kangbaek324/kkachi/internal/config"
	"github.com/kangbaek324/kkachi/internal/domain/currency"
	"github.com/kangbaek324/kkachi/internal/domain/user"
	"github.com/kangbaek324/kkachi/internal/domain/wallet"
)

func Register(r *gin.Engine, pool *pgxpool.Pool, cfg *config.Config) {
	r.NoRoute(func(c *gin.Context) {
		common.ApiResponse(c, http.StatusNotFound, false, nil, "Not found")
	})
	r.NoMethod(func(c *gin.Context) {
		common.ApiResponse(c, http.StatusMethodNotAllowed, false, nil, "Method not allowed")
	})

	r.GET("/health", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			common.ApiResponse(c, http.StatusServiceUnavailable, false, nil, "Database unavailable")
			return
		}
		common.ApiResponse(c, http.StatusOK, true, "ok")
	})

	v1 := r.Group("/api/v1")
	user.Register(v1, pool, cfg.JWTSecret)
	wallet.Register(v1, pool, cfg.JWTSecret)
	currency.Register(v1)
}
