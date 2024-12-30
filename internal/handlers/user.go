package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/response"
)

func (h *Handlers) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的输入",
		})
		return
	}

	user, err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		switch err {
		case errors.ErrDuplicateEntry:
			c.JSON(http.StatusConflict, gin.H{
				"code":    -1,
				"message": "记录已存在",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "服务器内部错误",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
	}))
}

func (h *Handlers) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的输入",
		})
		return
	}

	user, token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case errors.ErrNotFound, errors.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "未授权的访问",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "服务器内部错误",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	}))
}

func (h *Handlers) GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	user, err := h.userService.GetByID(userID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"code":    -1,
				"message": "记录不存在",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "服务器内部错误",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
	}))
}

func (h *Handlers) UpdateProfile(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的输入",
		})
		return
	}

	userID := c.GetInt64("user_id")
	updates := map[string]interface{}{
		"username": req.Username,
		"password": req.Password,
	}

	if err := h.userService.Update(userID, updates); err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"code":    -1,
				"message": "记录不存在",
			})
		case errors.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    -1,
				"message": "无效的输入",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "服务器内部错误",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}
