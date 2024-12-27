package api

import (
	"blockchain-shop/internal/api/middleware"
	"blockchain-shop/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *handlers.Handlers) {
	v1 := r.Group("/api/v1")
	{
		// 用户注册
		v1.POST("/users/register", h.Register)
		// 用户登录
		v1.POST("/users/login", h.Login)

		// 商品列表
		v1.GET("/products", h.ListProducts)
		// 商品详情
		v1.GET("/products/:id", h.GetProduct)
	}
}

func RegisterAuthRoutes(r *gin.Engine, h *handlers.Handlers, authMiddleware *middleware.JWTMiddleware) {
	v1 := r.Group("/api/v1")
	v1.Use(authMiddleware.MiddlewareFunc())
	{
		// 用户相关路由
		userGroup := v1.Group("/users")
		{
			userGroup.GET("/profile", h.GetProfile)
			userGroup.PUT("/profile", h.UpdateProfile)
		}

		// 商品相关路由
		productGroup := v1.Group("/products")
		{
			productGroup.POST("", h.CreateProduct)
			productGroup.PUT("/:id", h.UpdateProduct)
		}

		// 订单相关路由
		orderGroup := v1.Group("/orders")
		{
			orderGroup.POST("", h.CreateOrder)
			orderGroup.GET("/:id", h.GetOrder)
			orderGroup.GET("", h.ListOrders)
			orderGroup.GET("/:id/transaction", h.GetOrderTransaction)
		}
	}
}
