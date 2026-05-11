package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJSON[T any](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		ApiResponse(c, http.StatusBadRequest, false, nil, err.Error())
		return req, false
	}
	return req, true
}
