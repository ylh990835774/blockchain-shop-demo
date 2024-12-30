package middleware

import (
	"fmt"
	"net/http"

	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/logger"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 处理所有的错误响应
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) == 0 {
			return
		}

		// 获取最后一个错误
		err := c.Errors.Last()

		// 获取用户信息用于日志
		userID, exists := c.Get("user_id")
		userIDStr := "未登录"
		if exists {
			userIDStr = fmt.Sprintf("%v", userID)
		}

		// 记录错误日志
		logger.Error("API错误",
			logger.String("path", c.Request.URL.Path),
			logger.String("method", c.Request.Method),
			logger.String("user_id", userIDStr),
			logger.Err(err.Err),
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
				message = err.Error()
			case errors.ErrUnauthorized:
				statusCode = http.StatusUnauthorized
				message = err.Error()
			case errors.ErrForbidden:
				statusCode = http.StatusForbidden
				message = err.Error()
			case errors.ErrBadRequest:
				statusCode = http.StatusBadRequest
				message = err.Error()
			case errors.ErrInvalidInput:
				statusCode = http.StatusBadRequest
				message = err.Error()
			case errors.ErrDuplicateEntry:
				statusCode = http.StatusConflict
				message = err.Error()
			case errors.ErrNoFieldsToUpdate:
				statusCode = http.StatusBadRequest
				message = err.Error()
			default:
				statusCode = http.StatusInternalServerError
				message = "服务器内部错误"
			}
		}

		// 使用标准响应结构
		c.JSON(statusCode, response.Error(-1, message))
	}
}

// ResponseHandler 处理成功响应的中间件
func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果已经有错误处理了，就不需要处理成功响应
		if len(c.Errors) > 0 {
			return
		}

		// 如果响应已经被写入，就不需要处理
		if c.Writer.Written() {
			return
		}
	}
}
