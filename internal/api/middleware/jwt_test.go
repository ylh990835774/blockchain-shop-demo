package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ylh990835774/blockchain-shop-demo/internal/service"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

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

func TestJWTMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupAuth      func(*http.Request)
		setupMock      func(*MockJWTService)
		expectedStatus int
		expectedUserID interface{}
	}{
		{
			name: "valid token",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid.jwt.token")
			},
			setupMock: func(m *MockJWTService) {
				m.On("ParseToken", "valid.jwt.token").Return(int64(123), nil)
			},
			expectedStatus: http.StatusOK,
			expectedUserID: float64(123),
		},
		{
			name: "missing auth header",
			setupAuth: func(req *http.Request) {
				// 不设置Authorization头
			},
			setupMock:      func(m *MockJWTService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: nil,
		},
		{
			name: "invalid auth format",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "invalid_format")
			},
			setupMock:      func(m *MockJWTService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: nil,
		},
		{
			name: "invalid token",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid.token")
			},
			setupMock: func(m *MockJWTService) {
				m.On("ParseToken", "invalid.token").Return(int64(0), errors.ErrUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: nil,
		},
		{
			name: "malformed bearer token",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer")
			},
			setupMock:      func(m *MockJWTService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockJWTService)
			tt.setupMock(mockService)

			middleware := NewJWTMiddleware(mockService)
			router := gin.New()
			router.Use(middleware.MiddlewareFunc())
			router.GET("/test", func(c *gin.Context) {
				userID, exists := c.Get("user_id")
				if !exists {
					userID = nil
				}
				c.JSON(http.StatusOK, gin.H{"user_id": userID})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.setupAuth != nil {
				tt.setupAuth(req)
			}
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			if tt.expectedStatus == http.StatusOK {
				var result map[string]interface{}
				err := json.Unmarshal(resp.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserID, result["user_id"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestNewJWTMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		service     service.IJWTService
		shouldPanic bool
	}{
		{
			name:        "valid parameters",
			service:     new(MockJWTService),
			shouldPanic: false,
		},
		{
			name:        "nil service",
			service:     nil,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					NewJWTMiddleware(tt.service)
				})
			} else {
				assert.NotPanics(t, func() {
					middleware := NewJWTMiddleware(tt.service)
					assert.NotNil(t, middleware)
				})
			}
		})
	}
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockJWTService)
	mockService.On("ParseToken", "expired.token").Return(int64(0), errors.ErrUnauthorized)

	middleware := NewJWTMiddleware(mockService)
	router := gin.New()
	router.Use(middleware.MiddlewareFunc())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer expired.token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	mockService.AssertExpectations(t)
}
