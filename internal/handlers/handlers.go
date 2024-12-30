package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/service"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/response"
)

// Handlers 包含所有HTTP处理器
type Handlers struct {
	userService    service.IUserService
	jwtService     service.IJWTService
	productService service.IProductService
	orderService   service.IOrderService
}

// NewHandlers 创建一个新的Handlers实例
func NewHandlers(
	userService service.IUserService,
	jwtService service.IJWTService,
	productService service.IProductService,
	orderService service.IOrderService,
) *Handlers {
	return &Handlers{
		userService:    userService,
		jwtService:     jwtService,
		productService: productService,
		orderService:   orderService,
	}
}

// ListProducts 获取商品列表
func (h *Handlers) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	products, total, err := h.productService.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"total":    total,
		"products": products,
	}))
}

// GetProduct 获取商品详情
func (h *Handlers) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的商品ID"))
		return
	}

	product, err := h.productService.GetByID(id)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "商品不存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(product))
}

// CreateProduct 创建商品
func (h *Handlers) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}

	if err := h.productService.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		return
	}

	c.JSON(http.StatusCreated, response.Success(product))
}

// UpdateProduct 更新商品
func (h *Handlers) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的商品ID"))
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}

	if err := h.productService.Update(id, updates); err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "商品不存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// DeleteProduct 删除商品
func (h *Handlers) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的商品ID"))
		return
	}

	if err := h.productService.Delete(id); err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "商品不存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// CreateOrder 创建订单
func (h *Handlers) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的请求参数"))
		return
	}

	userID := c.GetInt64("userID")

	// 获取商品信息
	product, err := h.productService.GetByID(req.ProductID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "商品不存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	// 创建订单
	order := &model.Order{
		UserID:     userID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalPrice: product.Price * float64(req.Quantity),
		Status:     model.OrderStatusPending,
	}

	if err := h.orderService.Create(order); err != nil {
		switch err {
		case errors.ErrInsufficientStock:
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "库存不足"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(order))
}

// ListOrders 获取订单列表
func (h *Handlers) ListOrders(c *gin.Context) {
	userID := c.GetInt64("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.orderService.ListByUserID(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"total":  total,
		"orders": orders,
	}))
}

// GetOrder 获取订单详情
func (h *Handlers) GetOrder(c *gin.Context) {
	userID := c.GetInt64("userID")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的订单ID"))
		return
	}

	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "订单不存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	// 验证订单所属用户
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, "无权访问此订单"))
		return
	}

	c.JSON(http.StatusOK, response.Success(order))
}

// GetOrderTransaction 获取订单的区块链交易信息
func (h *Handlers) GetOrderTransaction(c *gin.Context) {
	userID := c.GetInt64("userID")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "无效的订单ID"))
		return
	}

	// 先获取订单信息，验证权限
	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "订单不存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	// 验证订单所属用户
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, "无权访问此订单"))
		return
	}

	// 获取交易信息
	transaction, err := h.orderService.GetTransaction(orderID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "交易信息不存在"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "服务器内部错误"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(transaction))
}
