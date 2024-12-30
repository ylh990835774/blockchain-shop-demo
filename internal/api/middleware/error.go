package middleware

import (
	"net/http"

	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 处理所有的错误响应
func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) == 0 {
			return
		}

		// 获取最后一个错误
		err := c.Errors.Last()

		// 记录错误日志
		log.Error("内部错误",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Error(err.Err),
		)

		// 根据错误类型返回不同的状态码和消息
		var statusCode int
		var message string

		if err.Type == gin.ErrorTypePublic {
			// 公开错误，直接返回给客户端
			statusCode = http.StatusInternalServerError
			if err.Meta != nil {
				if code, ok := err.Meta.(int); ok {
					statusCode = code
				}
			}
			message = err.Error()
		} else {
			// 私有错误，根据错误类型返回不同的状态码和消息
			switch err.Err {
			case errors.ErrNotFound:
				statusCode = http.StatusNotFound
				message = "记录不存在"
			case errors.ErrUnauthorized:
				statusCode = http.StatusUnauthorized
				message = "未授权的访问"
			case errors.ErrInvalidInput:
				statusCode = http.StatusBadRequest
				message = "无效的输入"
			case errors.ErrDuplicateEntry:
				statusCode = http.StatusConflict
				message = "记录已存在"
			default:
				statusCode = http.StatusInternalServerError
				message = "服务器内部错误"
			}
		}

		c.JSON(statusCode, gin.H{
			"code":    -1,
			"message": message,
		})
	}
}
