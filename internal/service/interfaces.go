package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
)

// IUserService 定义了用户服务的接口
type IUserService interface {
	Register(username, password string) (*model.User, error)
	Login(username, password string) (*model.User, error)
	GetByID(id int64) (*model.User, error)
	Update(id int64, updates map[string]interface{}) error
}

// IJWTService 定义了JWT服务的接口
type IJWTService interface {
	GenerateToken(userID int64) (string, error)
	ParseToken(token string) (int64, error)
}

// IProductService 定义了商品服务的接口
type IProductService interface {
	Create(product *model.Product) error
	Update(id int64, updates map[string]interface{}) error
	Delete(id int64) error
	Get(id int64) (*model.Product, error)
	List(page, pageSize int) ([]*model.Product, int64, error)
}

// IOrderService 定义了订单服务的接口
type IOrderService interface {
	Create(userID, productID int64, quantity int) (*model.Order, error)
	Get(id int64) (*model.Order, error)
	List(userID int64, page, pageSize int) ([]*model.Order, int64, error)
	GetTransaction(orderID int64) (*model.Transaction, error)
}
