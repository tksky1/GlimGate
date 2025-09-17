package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tksky1/glimgate/pkg/jwt"
	"github.com/tksky1/glimgate/pkg/response"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 解析token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Error(c, response.CodeInvalidToken)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)
		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			response.Forbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}
}