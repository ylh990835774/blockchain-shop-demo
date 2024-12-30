package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	customerrors "github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

// MockUserService 是用户服务的mock实现
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(username, password string) (*model.User, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Login(username, password string) (*model.User, string, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*model.User), args.String(1), args.Error(2)
}

func (m *MockUserService) GetByID(id int64) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Update(id int64, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockUserService) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// MockJWTService 是JWT服务的mock实现
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

// MockProductService 是商品服务的mock实现
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) Create(product *model.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) Update(id int64, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockProductService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductService) GetByID(id int64) (*model.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductService) List(page, pageSize int) ([]*model.Product, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

// MockOrderService 是订单服务的mock实现
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) Create(order *model.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderService) GetByID(id int64) (*model.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderService) ListByUserID(userID int64, page, pageSize int) ([]*model.Order, int64, error) {
	args := m.Called(userID, page, pageSize)
	return args.Get(0).([]*model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderService) GetTransaction(orderID int64) (*model.Transaction, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func TestHandlers_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*MockUserService, *MockJWTService)
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful_login",
			setupMock: func(m *MockUserService, j *MockJWTService) {
				m.On("Login", "testuser", "password123").Return(&model.User{
					ID:       1,
					Username: "testuser",
				}, "test.jwt.token", nil)
			},
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(200),
				"message": "success",
				"data": map[string]interface{}{
					"user": map[string]interface{}{
						"id":       float64(1),
						"username": "testuser",
					},
					"token": "test.jwt.token",
				},
			},
		},
		{
			name: "invalid_credentials",
			setupMock: func(m *MockUserService, j *MockJWTService) {
				m.On("Login", "testuser", "wrongpass").Return(nil, "", customerrors.ErrUnauthorized)
			},
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "wrongpass",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "未授权的访问",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(MockUserService)
			mockJWTService := new(MockJWTService)
			mockProductService := new(MockProductService)
			mockOrderService := new(MockOrderService)

			if tt.setupMock != nil {
				tt.setupMock(mockUserService, mockJWTService)
			}

			handlers := NewHandlers(mockUserService, mockJWTService, mockProductService, mockOrderService)
			router := gin.New()
			router.POST("/login", handlers.Login)

			requestJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestJSON))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			var actualBody map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualBody)

			mockUserService.AssertExpectations(t)
			mockJWTService.AssertExpectations(t)
		})
	}
}

func TestHandlers_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*MockUserService)
		userID         int64
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful_get_profile",
			setupMock: func(m *MockUserService) {
				m.On("GetByID", int64(1)).Return(&model.User{
					ID:       1,
					Username: "testuser",
				}, nil)
			},
			userID:         1,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(200),
				"message": "success",
				"data": map[string]interface{}{
					"id":       float64(1),
					"username": "testuser",
				},
			},
		},
		{
			name: "user_not_found",
			setupMock: func(m *MockUserService) {
				m.On("GetByID", int64(999)).Return(nil, customerrors.ErrNotFound)
			},
			userID:         999,
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "记录不存在",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(MockUserService)
			mockJWTService := new(MockJWTService)
			mockProductService := new(MockProductService)
			mockOrderService := new(MockOrderService)

			if tt.setupMock != nil {
				tt.setupMock(mockUserService)
			}

			handlers := NewHandlers(mockUserService, mockJWTService, mockProductService, mockOrderService)
			router := gin.New()
			router.GET("/profile", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				handlers.GetProfile(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			var actualBody map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualBody)

			mockUserService.AssertExpectations(t)
		})
	}
}

func TestHandlers_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*MockUserService)
		userID         int64
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful_update",
			setupMock: func(m *MockUserService) {
				m.On("Update", int64(1), mock.Anything).Return(nil)
			},
			userID: 1,
			requestBody: map[string]interface{}{
				"username": "newusername",
				"password": "newpassword",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(200),
				"message": "success",
				"data":    nil,
			},
		},
		{
			name: "invalid_request",
			setupMock: func(m *MockUserService) {
				m.On("Update", int64(1), mock.Anything).Return(customerrors.ErrInvalidInput)
			},
			userID: 1,
			requestBody: map[string]interface{}{
				"username": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "无效的输入",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(MockUserService)
			mockJWTService := new(MockJWTService)
			mockProductService := new(MockProductService)
			mockOrderService := new(MockOrderService)

			if tt.setupMock != nil {
				tt.setupMock(mockUserService)
			}

			handlers := NewHandlers(mockUserService, mockJWTService, mockProductService, mockOrderService)
			router := gin.New()
			router.PUT("/profile", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				handlers.UpdateProfile(c)
			})

			requestJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(requestJSON))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			var actualBody map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualBody)

			mockUserService.AssertExpectations(t)
		})
	}
}
