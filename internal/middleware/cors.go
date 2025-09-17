package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tksky1/glimgate/pkg/config"
)

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		corsConfig := config.AppConfig.CORS

		// 设置允许的源
		if len(corsConfig.AllowOrigins) > 0 {
			if corsConfig.AllowOrigins[0] == "*" {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				origin := c.GetHeader("Origin")
				for _, allowOrigin := range corsConfig.AllowOrigins {
					if origin == allowOrigin {
						c.Header("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}
		}

		// 设置允许的方法
		if len(corsConfig.AllowMethods) > 0 {
			methods := ""
			for i, method := range corsConfig.AllowMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			c.Header("Access-Control-Allow-Methods", methods)
		}

		// 设置允许的头部
		if len(corsConfig.AllowHeaders) > 0 {
			if corsConfig.AllowHeaders[0] == "*" {
				c.Header("Access-Control-Allow-Headers", "*")
			} else {
				headers := ""
				for i, header := range corsConfig.AllowHeaders {
					if i > 0 {
						headers += ", "
					}
					headers += header
				}
				c.Header("Access-Control-Allow-Headers", headers)
			}
		}

		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}