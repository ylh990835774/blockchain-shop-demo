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
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

// MockUserService 是 UserService 的 mock 实现
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

// MockJWTService 是 JWTService 的 mock 实现
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

func TestHandlers_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功注册", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockJWTService := new(MockJWTService)
		h := NewHandlers(mockUserService, mockJWTService, nil, nil)

		user := &model.User{
			ID:       1,
			Username: "testuser",
		}

		mockUserService.On("Register", "testuser", "password123").Return(user, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserService.AssertExpectations(t)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockJWTService := new(MockJWTService)
		h := NewHandlers(mockUserService, mockJWTService, nil, nil)

		mockUserService.On("Register", "existinguser", "password123").Return(nil, errors.ErrDuplicateEntry)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := map[string]string{
			"username": "existinguser",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockUserService.AssertExpectations(t)
	})
}

func TestHandlers_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUserService, *MockJWTService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful login",
			requestBody: LoginRequest{
				Username: "testuser",
				Password: "testpass",
			},
			setupMock: func(m *MockUserService, j *MockJWTService) {
				m.On("Login", "testuser", "testpass").Return(&model.User{
					ID:       1,
					Username: "testuser",
				}, "test.jwt.token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    200,
				"message": "success",
				"data": map[string]interface{}{
					"token": "test.jwt.token",
					"user": map[string]interface{}{
						"id":       float64(1),
						"username": "testuser",
					},
				},
			},
		},
		{
			name: "invalid credentials",
			requestBody: LoginRequest{
				Username: "wronguser",
				Password: "wrongpass",
			},
			setupMock: func(m *MockUserService, j *MockJWTService) {
				m.On("Login", "wronguser", "wrongpass").Return(nil, "", errors.ErrUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"code":    401,
				"message": "用户名或密码错误",
				"data":    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(MockUserService)
			mockJWTService := new(MockJWTService)
			tt.setupMock(mockUserService, mockJWTService)

			h := NewHandlers(mockUserService, mockJWTService, nil, nil)

			router := gin.New()
			router.POST("/login", h.Login)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			var actualBody map[string]interface{}
			json.Unmarshal(resp.Body.Bytes(), &actualBody)
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
		userID         int64
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful get profile",
			userID: 1,
			setupMock: func(m *MockUserService) {
				m.On("GetByID", int64(1)).Return(&model.User{
					ID:       1,
					Username: "testuser",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(0),
				"message": "success",
				"data": map[string]interface{}{
					"id":       float64(1),
					"username": "testuser",
				},
			},
		},
		{
			name:   "user not found",
			userID: 999,
			setupMock: func(m *MockUserService) {
				m.On("GetByID", int64(999)).Return(nil, errors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "用户不存在",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(MockUserService)
			tt.setupMock(mockUserService)

			h := &Handlers{
				userService: mockUserService,
			}

			router := gin.New()
			router.GET("/profile", func(c *gin.Context) {
				c.Set("userID", tt.userID)
				h.GetProfile(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			var actualBody map[string]interface{}
			json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.Equal(t, tt.expectedBody, actualBody)

			mockUserService.AssertExpectations(t)
		})
	}
}

func TestHandlers_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         int64
		requestBody    interface{}
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful update",
			userID: 1,
			requestBody: UpdateProfileRequest{
				Phone:   "1234567890",
				Address: "test address",
			},
			setupMock: func(m *MockUserService) {
				m.On("Update", int64(1), map[string]interface{}{
					"phone":   "1234567890",
					"address": "test address",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(0),
				"message": "success",
			},
		},
		{
			name:   "invalid request",
			userID: 1,
			requestBody: gin.H{
				"phone":   123, // 错误的类型
				"address": 456, // 错误的类型
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "无效的请求参数",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(MockUserService)
			tt.setupMock(mockUserService)

			h := &Handlers{
				userService: mockUserService,
			}

			router := gin.New()
			router.PUT("/profile", func(c *gin.Context) {
				c.Set("userID", tt.userID)
				h.UpdateProfile(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			var actualBody map[string]interface{}
			json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.Equal(t, tt.expectedBody, actualBody)

			mockUserService.AssertExpectations(t)
		})
	}
}
