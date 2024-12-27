package handlers

import (
	"blockchain-shop/internal/blockchain"
	"blockchain-shop/internal/service"
)

type Handlers struct {
	userService       *service.UserService
	productService    *service.ProductService
	orderService      *service.OrderService
	blockchainService blockchain.Service
	jwtService        *service.JWTService
}

func NewHandlers(userService *service.UserService, productService *service.ProductService, orderService *service.OrderService, blockchainService blockchain.Service, jwtService *service.JWTService) *Handlers {
	return &Handlers{
		userService:       userService,
		productService:    productService,
		orderService:      orderService,
		blockchainService: blockchainService,
		jwtService:        jwtService,
	}
}
