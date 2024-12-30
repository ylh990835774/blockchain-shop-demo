package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/service"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
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
		handleError(c, err, "获取商品列表")
		return
	}

	handleSuccess(c, gin.H{
		"total":    total,
		"products": products,
	}, "获取商品列表")
}

// GetProduct 获取商品详情
func (h *Handlers) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		handleError(c, errors.ErrInvalidInput, "获取商品详情-参数验证")
		return
	}

	product, err := h.productService.GetByID(id)
	if err != nil {
		handleError(c, err, "获取商品详情")
		return
	}

	handleSuccess(c, product, "获取商品详情")
}

// CreateProduct 创建商品
func (h *Handlers) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		handleError(c, errors.ErrInvalidInput, "创建商品-参数验证")
		return
	}

	if err := h.productService.Create(&product); err != nil {
		handleError(c, err, "创建商品")
		return
	}

	handleSuccess(c, product, "创建商品")
}

// UpdateProduct 更新商品
func (h *Handlers) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		handleError(c, errors.ErrInvalidInput, "更新商品-参数验证")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		handleError(c, errors.ErrInvalidInput, "更新商品-参数验证")
		return
	}

	if err := h.productService.Update(id, updates); err != nil {
		handleError(c, err, "更新商品")
		return
	}

	handleSuccess(c, nil, "更新商品")
}

// DeleteProduct 删除商品
func (h *Handlers) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		handleError(c, errors.ErrInvalidInput, "删除商品-参数验证")
		return
	}

	if err := h.productService.Delete(id); err != nil {
		handleError(c, err, "删除商品")
		return
	}

	handleSuccess(c, nil, "删除商品")
}

// CreateOrder 创建订单
func (h *Handlers) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.ErrInvalidInput, "创建订单-参数验证")
		return
	}

	userID := c.GetInt64("user_id")

	product, err := h.productService.GetByID(req.ProductID)
	if err != nil {
		handleError(c, err, "创建订单-获取商品信息")
		return
	}

	order := &model.Order{
		UserID:     userID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalPrice: product.Price * float64(req.Quantity),
		Status:     model.OrderStatusPending,
	}

	if err := h.orderService.Create(order); err != nil {
		handleError(c, err, "创建订单")
		return
	}

	handleSuccess(c, order, "创建订单")
}

// ListOrders 获取订单列表
func (h *Handlers) ListOrders(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.orderService.ListByUserID(userID, page, pageSize)
	if err != nil {
		handleError(c, err, "获取订单列表")
		return
	}

	handleSuccess(c, gin.H{
		"total":  total,
		"orders": orders,
	}, "获取订单列表")
}

// GetOrder 获取订单详情
func (h *Handlers) GetOrder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		handleError(c, errors.ErrInvalidInput, "获取订单详情-参数验证")
		return
	}

	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		handleError(c, err, "获取订单详情")
		return
	}

	// 验证订单所属用户
	if order.UserID != userID {
		handleError(c, errors.ErrUnauthorized, "获取订单详情-权限验证")
		return
	}

	handleSuccess(c, order, "获取订单详情")
}

// GetOrderTransaction 获取订单的区块链交易信息
func (h *Handlers) GetOrderTransaction(c *gin.Context) {
	userID := c.GetInt64("user_id")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		handleError(c, errors.ErrInvalidInput, "获取订单交易信息-参数验证")
		return
	}

	// 先获取订单信息，验证权限
	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		handleError(c, err, "获取订单交易信息-订单验证")
		return
	}

	// 验证订单所属用户
	if order.UserID != userID {
		handleError(c, errors.ErrUnauthorized, "获取订单交易信息-权限验证")
		return
	}

	// 获取交易信息
	transaction, err := h.orderService.GetTransaction(orderID)
	if err != nil {
		handleError(c, err, "获取订单交易信息")
		return
	}

	handleSuccess(c, transaction, "获取订单交易信息")
}
