package common

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ApiResponse(c *gin.Context, code int, success bool, data any, message ...string) {
	m := "요청이 성공적으로 처리되었습니다."
	if len(message) > 0 {
		m = message[0]
	}

	c.JSON(code, Response{
		Code:    code,
		Success: success,
		Message: m,
		Data:    data,
	})
}
