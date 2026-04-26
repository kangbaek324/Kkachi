package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	_ = api
	// api.Use(middleware.Auth(jwtSecret))
	// Register routes here
}
