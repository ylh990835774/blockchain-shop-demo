package middleware

import (
	"net/http"

	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/logger"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/response"

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
				message = "资源未找到"
			case errors.ErrUnauthorized:
				statusCode = http.StatusUnauthorized
				message = "未授权访问"
			case errors.ErrForbidden:
				statusCode = http.StatusForbidden
				message = "禁止访问"
			case errors.ErrBadRequest:
				statusCode = http.StatusBadRequest
				message = "请求参数错误"
			default:
				statusCode = http.StatusInternalServerError
				message = "服务器内部错误"
			}
		}

		// 使用标准响应结构
		c.JSON(statusCode, response.Error(statusCode, message))
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

		// 获取原始数据
		if data, exists := c.Get("data"); exists {
			// 使用标准响应结构
			c.JSON(http.StatusOK, response.Success(data))
			return
		}

		// 如果没有设置数据，直接返回，让原始的处理器处理响应
		return
	}
}
