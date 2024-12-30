package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/internal/api/middleware"
	"github.com/ylh990835774/blockchain-shop-demo/internal/handlers"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/logger"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine, h *handlers.Handlers, log *logger.Logger, jwtMiddleware *middleware.JWTMiddleware) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 用户公开接口
		users := v1.Group("/users")
		{
			users.POST("/register", h.Register)
			users.POST("/login", h.Login)
		}

		// 公开的商品列表和详情接口
		v1.GET("/products", h.ListProducts)
		v1.GET("/products/:id", h.GetProduct)

		// 需要认证的接口
		auth := v1.Group("")
		auth.Use(jwtMiddleware.MiddlewareFunc())
		{
			// 用户认证接口
			authUsers := auth.Group("/users")
			{
				authUsers.GET("/profile", h.GetProfile)
				authUsers.PUT("/profile", h.UpdateProfile)
			}

			// 商品相关接口
			products := auth.Group("/products")
			{
				products.POST("", h.CreateProduct)
				products.PUT("/:id", h.UpdateProduct)
				products.DELETE("/:id", h.DeleteProduct)
			}

			// 订单相关接口
			orders := auth.Group("/orders")
			{
				orders.POST("", h.CreateOrder)
				orders.GET("", h.ListOrders)
				orders.GET("/:id", h.GetOrder)
				orders.GET("/:id/transaction", h.GetOrderTransaction)
			}
		}
	}
}
