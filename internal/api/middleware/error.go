package middleware

import (
	"net/http"

	"blockchain-shop/pkg/errors"
	"blockchain-shop/pkg/logger"

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

		// 根据错误类型处理
		if err.Type == gin.ErrorTypePublic {
			// 公开错误，直接返回给客户端
			statusCode := http.StatusInternalServerError
			if err.Meta != nil {
				if code, ok := err.Meta.(int); ok {
					statusCode = code
				}
			}

			c.JSON(statusCode, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
			return
		}

		// 私有错误，记录日志并返回通用错误信息
		log.Error("内部错误",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Error(err.Err),
		)

		// 根据错误类型返回不同的状态码
		var statusCode int
		switch err.Err {
		case errors.ErrNotFound:
			statusCode = http.StatusNotFound
		case errors.ErrUnauthorized:
			statusCode = http.StatusUnauthorized
		case errors.ErrInvalidInput:
			statusCode = http.StatusBadRequest
		case errors.ErrDuplicateEntry:
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{
			"code":    -1,
			"message": "服务器内部错误",
		})
	}
}
