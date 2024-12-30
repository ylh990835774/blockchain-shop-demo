package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id int64) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(id int64, updates interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByUsername(username string) bool {
	args := m.Called(username)
	return args.Bool(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID int64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ParseToken(token string) (int64, error) {
	args := m.Called(token)
	return args.Get(0).(int64), args.Error(1)
}

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	service := NewUserService(mockRepo, mockJWT)

	t.Run("成功注册", func(t *testing.T) {
		username := "testuser"
		password := "password123"

		mockRepo.On("GetByUsername", username).Return(nil, errors.ErrNotFound)
		mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

		user, err := service.Register(username, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, username, user.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		username := "existinguser"
		password := "password123"

		existingUser := &model.User{
			ID:       1,
			Username: username,
			Password: "hashedpassword",
		}

		mockRepo.On("GetByUsername", username).Return(existingUser, nil)

		user, err := service.Register(username, password)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrDuplicateEntry, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}
