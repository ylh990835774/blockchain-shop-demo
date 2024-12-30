package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
)

// IUserService 用户服务接口
type IUserService interface {
	Register(username, password string) (*model.User, error)
	Login(username, password string) (*model.User, string, error) // 返回用户信息和token
	GetByID(id int64) (*model.User, error)
	Update(id int64, updates map[string]interface{}) error
	GetByUsername(username string) (*model.User, error)
}

// IJWTService JWT服务接口
type IJWTService interface {
	GenerateToken(userID int64) (string, error)
	ParseToken(token string) (int64, error)
}

// IProductService 商品服务接口
type IProductService interface {
	Create(product *model.Product) error
	Update(id int64, updates map[string]interface{}) error
	Delete(id int64) error
	GetByID(id int64) (*model.Product, error)
	List(page, pageSize int) ([]*model.Product, int64, error)
}

// IOrderService 订单服务接口
type IOrderService interface {
	Create(order *model.Order) error
	GetByID(id int64) (*model.Order, error)
	ListByUserID(userID int64, page, pageSize int) ([]*model.Order, int64, error)
	GetTransaction(orderID int64) (*model.Transaction, error)
}
