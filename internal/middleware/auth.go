package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetBearerToken은 환경변수에서 Bearer token을 가져옵니다
func GetBearerToken() (string, error) {
	token := os.Getenv("BEARER_TOKEN")
	if token == "" {
		return "", errors.New("BEARER_TOKEN 환경변수가 비어있습니다")
	}
	return token, nil
}

// BearerTokenAuth 미들웨어는 Bearer token을 검증합니다
func BearerTokenAuth() gin.HandlerFunc {
	expectedToken, err := GetBearerToken()
	if err != nil {
		log.Fatalf("토큰을 받아오는데 실패했습니다: %v", err)
	}
	
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

		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		// 토큰 검증
		if token != expectedToken {
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

