package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

// UserService 实现了IUserService接口
type userService struct {
	repo repository.UserRepository
}

// 确保userService实现了IUserService接口
var _ IUserService = (*userService)(nil)

func NewUserService(repo repository.UserRepository) IUserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Create(user *model.User) error {
	if user == nil {
		return errors.ErrInvalidInput
	}

	if s.ExistsByUsername(user.Username) {
		return errors.ErrDuplicateEntry
	}

	return s.repo.Create(user)
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

func (s *userService) Authenticate(username, password string) (*model.User, error) {
	user, err := s.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, errors.ErrUnauthorized
	}

	return user, nil
}

func (s *userService) ExistsByUsername(username string) bool {
	return s.repo.ExistsByUsername(username)
}

func (s *userService) Update(id int64, updates map[string]interface{}) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	if _, err := s.GetByID(id); err != nil {
		if err == errors.ErrNotFound {
			return errors.ErrNotFound
		}
		return err
	}

	return s.repo.Update(id, updates)
}

func (s *userService) Register(username, password string) (*model.User, error) {
	// 检查用户是否已存在
	if s.ExistsByUsername(username) {
		return nil, errors.ErrDuplicateEntry
	}

	// 创建用户
	user := &model.User{
		Username: username,
	}
	user.SetPassword(password)

	if err := s.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(username, password string) (*model.User, error) {
	user, err := s.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, errors.ErrUnauthorized
	}

	return user, nil
}

func (s *userService) GetByUsername(username string) (*model.User, error) {
	if username == "" {
		return nil, errors.ErrInvalidInput
	}

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}
