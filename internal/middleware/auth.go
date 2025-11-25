package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const BearerToken = "gistsikdang"

// BearerTokenAuth 미들웨어는 Bearer token을 검증합니다
func BearerTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header is required",
			})
			c.Abort()
			return
		}

		// "Bearer " 접두사 제거
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			// Bearer 접두사가 없었음
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid authorization format. Use 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		// 토큰 검증
		if token != BearerToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

