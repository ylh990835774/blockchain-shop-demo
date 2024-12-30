package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/logger"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/response"
)

// handleError 统一处理错误，记录日志并返回适当的响应
func handleError(c *gin.Context, err error, operation string) {
	// 获取请求信息用于日志
	path := c.Request.URL.Path
	method := c.Request.Method
	userID, exists := c.Get("user_id")
	userIDStr := "未登录"
	if exists {
		userIDStr = fmt.Sprintf("%v", userID)
	}

	// 记录错误日志
	logger.Error("API错误",
		logger.String("path", path),
		logger.String("method", method),
		logger.String("user_id", userIDStr),
		logger.String("operation", operation),
		logger.Err(err),
	)

	// 根据错误类型返回不同的响应
	switch err {
	case errors.ErrNotFound:
		c.JSON(http.StatusNotFound, response.Error(-1, err.Error()))
	case errors.ErrUnauthorized:
		c.JSON(http.StatusUnauthorized, response.Error(-1, err.Error()))
	case errors.ErrInvalidInput:
		c.JSON(http.StatusBadRequest, response.Error(-1, err.Error()))
	case errors.ErrDuplicateEntry:
		c.JSON(http.StatusConflict, response.Error(-1, err.Error()))
	case errors.ErrNoFieldsToUpdate:
		c.JSON(http.StatusBadRequest, response.Error(-1, err.Error()))
	default:
		// 对于未知错误，返回500但记录详细日志
		logger.Error("未处理的错误",
			logger.String("path", path),
			logger.String("method", method),
			logger.String("user_id", userIDStr),
			logger.String("operation", operation),
			logger.Err(err),
		)
		c.JSON(http.StatusInternalServerError, response.Error(-1, "服务器内部错误"))
	}
}

// handleSuccess 统一处理成功响应，可选是否记录日志
func handleSuccess(c *gin.Context, data interface{}, operation string) {
	// 可以选择记录成功操作的日志
	if operation != "" {
		path := c.Request.URL.Path
		method := c.Request.Method
		userID, exists := c.Get("user_id")
		userIDStr := "未登录"
		if exists {
			userIDStr = fmt.Sprintf("%v", userID)
		}

		logger.Info("API操作成功",
			logger.String("path", path),
			logger.String("method", method),
			logger.String("user_id", userIDStr),
			logger.String("operation", operation),
		)
	}

	c.JSON(http.StatusOK, response.Success(data))
}
