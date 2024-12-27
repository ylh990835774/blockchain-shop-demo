package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/blockchain"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
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
