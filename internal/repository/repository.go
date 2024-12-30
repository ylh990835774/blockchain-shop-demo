package repository

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
)

// UserRepository 定义用户仓库接口
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id int64) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	ExistsByUsername(username string) bool
	Update(id int64, updates interface{}) error
}
