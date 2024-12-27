package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("invalid request"),
			Meta: http.StatusBadRequest,
		})
		return
	}

	user, err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": user})
}

func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("invalid request"),
			Meta: http.StatusBadRequest,
		})
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"token": token,
		"user":  user,
	}})
}

func (h *Handlers) GetProfile(c *gin.Context) {
	userID := c.GetInt64("userID")
	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("user not found"),
			Meta: http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": user})
}

func (h *Handlers) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("userID")
	var updateReq struct {
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  err,
			Meta: http.StatusBadRequest,
		})
		return
	}

	updates := map[string]interface{}{
		"phone":   updateReq.Phone,
		"address": updateReq.Address,
	}

	if err := h.userService.Update(userID, updates); err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
