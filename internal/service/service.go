package service

import (
	"blockchain-shop/internal/blockchain"
	"blockchain-shop/internal/repository/mysql"
)

type Services struct {
	User    UserService
	Product ProductService
	Order   OrderService
}

func NewServices(
	userRepo *mysql.UserRepository,
	productRepo *mysql.ProductRepository,
	orderRepo *mysql.OrderRepository,
	blockchainSvc blockchain.Service,
) *Services {
	productService := NewProductService(productRepo)
	return &Services{
		User:    *NewUserService(userRepo),
		Product: *productService,
		Order:   *NewOrderService(orderRepo, productService, blockchainSvc),
	}
}
