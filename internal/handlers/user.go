package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/response"
)

func (h *Handlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}

	user, err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		switch err {
		case errors.ErrDuplicateEntry:
			c.JSON(http.StatusConflict, response.Error(http.StatusConflict, "用户名已存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
	}))
}

func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}

	user, token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case errors.ErrNotFound, errors.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "用户名或密码错误"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
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
	userID := c.GetInt64("userID")
	user, err := h.userService.GetByID(userID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "用户不存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
	}))
}

func (h *Handlers) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的请求参数"))
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
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "用户不存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}
