package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

// MockUserRepository 是UserRepository的mock实现
type MockUserRepository struct {
	mock.Mock
}

// 确保MockUserRepository实现了repository.UserRepository接口
var _ repository.UserRepository = (*MockUserRepository)(nil)

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

func (m *MockUserRepository) ExistsByUsername(username string) bool {
	args := m.Called(username)
	return args.Bool(0)
}

func (m *MockUserRepository) Update(id int64, updates interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	tests := []struct {
		name         string
		username     string
		password     string
		mockBehavior func()
		expectError  error
		expectUser   bool
	}{
		{
			name:     "successful registration",
			username: "testuser",
			password: "password123",
			mockBehavior: func() {
				mockRepo.On("ExistsByUsername", "testuser").Return(false)
				mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
			},
			expectError: nil,
			expectUser:  true,
		},
		{
			name:     "duplicate username",
			username: "existinguser",
			password: "password123",
			mockBehavior: func() {
				mockRepo.On("ExistsByUsername", "existinguser").Return(true)
			},
			expectError: errors.ErrDuplicateEntry,
			expectUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockBehavior()

			user, err := service.Register(tt.username, tt.password)

			if tt.expectError != nil {
				assert.Equal(t, tt.expectError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
				assert.True(t, user.CheckPassword(tt.password))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	validUser := &model.User{
		Username: "testuser",
	}
	validUser.SetPassword("password123")

	tests := []struct {
		name         string
		username     string
		password     string
		mockBehavior func()
		expectError  error
		expectUser   bool
	}{
		{
			name:     "successful login",
			username: "testuser",
			password: "password123",
			mockBehavior: func() {
				mockRepo.On("GetByUsername", "testuser").Return(validUser, nil)
			},
			expectError: nil,
			expectUser:  true,
		},
		{
			name:     "user not found",
			username: "nonexistent",
			password: "password123",
			mockBehavior: func() {
				mockRepo.On("GetByUsername", "nonexistent").Return(nil, mysql.ErrNotFound)
			},
			expectError: errors.ErrNotFound,
			expectUser:  false,
		},
		{
			name:     "wrong password",
			username: "testuser",
			password: "wrongpassword",
			mockBehavior: func() {
				mockRepo.On("GetByUsername", "testuser").Return(validUser, nil)
			},
			expectError: errors.ErrUnauthorized,
			expectUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockBehavior()

			user, err := service.Login(tt.username, tt.password)

			if tt.expectError != nil {
				assert.Equal(t, tt.expectError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &model.User{
		ID:       1,
		Username: "testuser",
	}

	tests := []struct {
		name         string
		id           int64
		mockBehavior func()
		expectError  error
		expectUser   bool
	}{
		{
			name: "successful get",
			id:   1,
			mockBehavior: func() {
				mockRepo.On("GetByID", int64(1)).Return(user, nil)
			},
			expectError: nil,
			expectUser:  true,
		},
		{
			name:         "invalid id",
			id:           0,
			mockBehavior: func() {},
			expectError:  errors.ErrInvalidInput,
			expectUser:   false,
		},
		{
			name: "user not found",
			id:   999,
			mockBehavior: func() {
				mockRepo.On("GetByID", int64(999)).Return(nil, mysql.ErrNotFound)
			},
			expectError: errors.ErrNotFound,
			expectUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockBehavior()

			user, err := service.GetByID(tt.id)

			if tt.expectError != nil {
				assert.Equal(t, tt.expectError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &model.User{
		ID:       1,
		Username: "testuser",
	}

	tests := []struct {
		name         string
		id           int64
		updates      map[string]interface{}
		mockBehavior func()
		expectError  error
	}{
		{
			name: "successful update",
			id:   1,
			updates: map[string]interface{}{
				"username": "newusername",
			},
			mockBehavior: func() {
				mockRepo.On("GetByID", int64(1)).Return(user, nil)
				mockRepo.On("Update", int64(1), mock.AnythingOfType("map[string]interface {}")).Return(nil)
			},
			expectError: nil,
		},
		{
			name: "invalid id",
			id:   0,
			updates: map[string]interface{}{
				"username": "newusername",
			},
			mockBehavior: func() {},
			expectError:  errors.ErrInvalidInput,
		},
		{
			name: "user not found",
			id:   999,
			updates: map[string]interface{}{
				"username": "newusername",
			},
			mockBehavior: func() {
				mockRepo.On("GetByID", int64(999)).Return(nil, mysql.ErrNotFound)
			},
			expectError: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockBehavior()

			err := service.Update(tt.id, tt.updates)

			if tt.expectError != nil {
				assert.Equal(t, tt.expectError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
