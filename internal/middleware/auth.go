package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kangbaek324/kkachi/internal/common"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "missing or invalid authorization header",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "invalid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "invalid token claims",
			})
			return
		}

		sub, err := claims.GetSubject()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Response{
				Code:    http.StatusInternalServerError,
				Success: false,
				Message: "internal server error",
			})
			return
		}

		userID, err := strconv.ParseInt(sub, 10, 64)
		c.Set("userId", userID)
		c.Next()
	}
}
