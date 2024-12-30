package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/logger"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 获取用户信息
		userID, exists := c.Get("user_id")
		userIDStr := "未登录"
		if exists {
			userIDStr = fmt.Sprintf("%v", userID)
		}

		c.Next()

		cost := time.Since(start)
		logger.Info("请求日志",
			logger.String("status", fmt.Sprintf("%d", c.Writer.Status())),
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", query),
			logger.String("ip", c.ClientIP()),
			logger.String("user-agent", c.Request.UserAgent()),
			logger.String("cost", cost.String()),
			logger.String("user_id", userIDStr),
		)
	}
}
