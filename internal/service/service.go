package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
)

// UserService 定义用户服务接口
type UserService interface {
	Register(username, password string) (*model.User, error)
	Login(username, password string) (*model.User, error)
	GetByID(id int64) (*model.User, error)
	Update(id int64, updates map[string]interface{}) error
}

// JWTService 定义JWT服务接口
type JWTService interface {
	GenerateToken(userID int64) (string, error)
}
