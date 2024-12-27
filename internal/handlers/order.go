package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"blockchain-shop/internal/model"
	pkgerrors "blockchain-shop/pkg/errors"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Validation error, return directly
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  err,
			Meta: http.StatusBadRequest,
		})
		return
	}

	order := &model.Order{
		UserID:    c.GetInt64("userID"),
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := h.orderService.Create(order); err != nil {
		// Internal error, add error, middleware will log and return generic response
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}

func (h *Handlers) ListOrders(c *gin.Context) {
	userID := c.GetInt64("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.orderService.ListByUserID(userID, page, pageSize)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"orders":   orders,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}})
}

func (h *Handlers) GetOrder(c *gin.Context) {
	userID := c.GetInt64("userID")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("invalid order id"),
			Meta: http.StatusBadRequest,
		})
		return
	}

	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("order not found"),
			Meta: http.StatusNotFound,
		})
		return
	}

	// 验证订单所属
	if order.UserID != userID {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("forbidden"),
			Meta: http.StatusForbidden,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}

// GetOrderTransaction 获取订单的区块链交易信息
func (h *Handlers) GetOrderTransaction(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	// 验证用户是否有权限查看���订单
	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		if err == pkgerrors.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 验证用户是否有权限查看该订单
	userID := c.GetInt64("userID")
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// 获取区块链交易信息
	txData, err := h.orderService.GetOrderTransaction(orderID)
	if err != nil {
		if err == pkgerrors.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var txInfo map[string]interface{}
	if err := json.Unmarshal(txData, &txInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse transaction data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id": orderID,
		"tx_hash":  order.TxHash,
		"tx_data":  txInfo,
	})
}
