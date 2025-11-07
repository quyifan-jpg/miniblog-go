package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"miniblog/internal/pkg/jwt"
	"miniblog/internal/pkg/response"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			response.Unauthorized(c, "token is empty")
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetUserID 从上下文中获取用户 ID
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

