package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"blockchain-shop/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  err,
			Meta: http.StatusBadRequest,
		})
		return
	}

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := h.productService.Create(product); err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}

func (h *Handlers) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	products, total, err := h.productService.List(page, pageSize)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}})
}

func (h *Handlers) GetProduct(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("invalid product id"),
			Meta: http.StatusBadRequest,
		})
		return
	}

	product, err := h.productService.GetByID(productID)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("product not found"),
			Meta: http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}

func (h *Handlers) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("invalid product id"),
			Meta: http.StatusBadRequest,
		})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  err,
			Meta: http.StatusBadRequest,
		})
		return
	}

	product := &model.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := h.productService.Update(product); err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}

func (h *Handlers) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New("invalid product id"),
			Meta: http.StatusBadRequest,
		})
		return
	}

	if err := h.productService.Delete(id); err != nil {
		c.Error(gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
