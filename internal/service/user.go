package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

// UserService 实现了IUserService接口
type userService struct {
	repo       repository.UserRepository
	jwtService IJWTService
}

// 确保userService实现了IUserService接口
var _ IUserService = (*userService)(nil)

func NewUserService(repo repository.UserRepository, jwtService IJWTService) IUserService {
	return &userService{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *userService) Register(username, password string) (*model.User, error) {
	// 检查用户名是否已存在
	if s.ExistsByUsername(username) {
		return nil, errors.ErrDuplicateEntry
	}

	// 创建新用户
	user := &model.User{
		Username: username,
	}

	// 设置加密密码
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(username, password string) (*model.User, string, error) {
	// 根据用户名获取用户
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, "", errors.ErrNotFound
		}
		return nil, "", err
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return nil, "", errors.ErrUnauthorized
	}

	// 生成JWT令牌
	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *userService) GetByID(id int64) (*model.User, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(id int64, updates map[string]interface{}) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	if err := s.repo.Update(id, updates); err != nil {
		if err == mysql.ErrNotFound {
			return errors.ErrNotFound
		}
		return err
	}

	return nil
}

func (s *userService) GetByUsername(username string) (*model.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) ExistsByUsername(username string) bool {
	_, err := s.repo.GetByUsername(username)
	return err == nil
}
