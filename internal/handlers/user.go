package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

func (h *Handlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的请求参数",
		})
		return
	}

	user, err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		switch err {
		case errors.ErrDuplicateEntry:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"code":    -1,
				"message": "用户名已存在",
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "服务器内部错误",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的请求参数",
		})
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case errors.ErrUnauthorized:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "用户名或密码错误",
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "服务器内部错误",
			})
		}
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
			},
		},
	})
}

func (h *Handlers) GetProfile(c *gin.Context) {
	userID := c.GetInt64("userID")
	user, err := h.userService.GetByID(userID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    -1,
				"message": "用户不存在",
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "服务器内部错误",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *Handlers) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的请求参数",
		})
		return
	}

	userID := c.GetInt64("userID")
	updates := map[string]interface{}{
		"phone":   req.Phone,
		"address": req.Address,
	}

	if err := h.userService.Update(userID, updates); err != nil {
		switch err {
		case errors.ErrNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    -1,
				"message": "用户不存在",
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "服务器内部错误",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
