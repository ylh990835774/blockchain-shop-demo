package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	customerrors "github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupRouter    func(*gin.Engine, *logger.Logger)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "public error with status code",
			setupRouter: func(r *gin.Engine, log *logger.Logger) {
				r.Use(ErrorHandler(log))
				r.GET("/test", func(c *gin.Context) {
					c.Error(&gin.Error{
						Type: gin.ErrorTypePublic,
						Err:  customerrors.ErrInvalidInput,
						Meta: http.StatusBadRequest,
					})
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "无效的输入",
			},
		},
		{
			name: "not found error",
			setupRouter: func(r *gin.Engine, log *logger.Logger) {
				r.Use(ErrorHandler(log))
				r.GET("/test", func(c *gin.Context) {
					c.Error(&gin.Error{
						Type: gin.ErrorTypePrivate,
						Err:  customerrors.ErrNotFound,
					})
				})
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "记录不存在",
			},
		},
		{
			name: "unauthorized error",
			setupRouter: func(r *gin.Engine, log *logger.Logger) {
				r.Use(ErrorHandler(log))
				r.GET("/test", func(c *gin.Context) {
					c.Error(&gin.Error{
						Type: gin.ErrorTypePrivate,
						Err:  customerrors.ErrUnauthorized,
					})
				})
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "未授权的访问",
			},
		},
		{
			name: "invalid input error",
			setupRouter: func(r *gin.Engine, log *logger.Logger) {
				r.Use(ErrorHandler(log))
				r.GET("/test", func(c *gin.Context) {
					c.Error(&gin.Error{
						Type: gin.ErrorTypePrivate,
						Err:  customerrors.ErrInvalidInput,
					})
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "无效的输入",
			},
		},
		{
			name: "duplicate entry error",
			setupRouter: func(r *gin.Engine, log *logger.Logger) {
				r.Use(ErrorHandler(log))
				r.GET("/test", func(c *gin.Context) {
					c.Error(&gin.Error{
						Type: gin.ErrorTypePrivate,
						Err:  customerrors.ErrDuplicateEntry,
					})
				})
			},
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "记录已存在",
			},
		},
		{
			name: "unknown error",
			setupRouter: func(r *gin.Engine, log *logger.Logger) {
				r.Use(ErrorHandler(log))
				r.GET("/test", func(c *gin.Context) {
					c.Error(&gin.Error{
						Type: gin.ErrorTypePrivate,
						Err:  errors.New("unknown error"),
					})
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"code":    float64(-1),
				"message": "服务器内部错误",
			},
		},
		{
			name: "no error",
			setupRouter: func(r *gin.Engine, log *logger.Logger) {
				r.Use(ErrorHandler(log))
				r.GET("/test", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code":    0,
						"message": "success",
					})
				})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(0),
				"message": "success",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试日志记录器
			testLogger := logger.New(zaptest.NewLogger(t))

			// 创建路由器
			router := gin.New()
			tt.setupRouter(router, testLogger)

			// 创建测试请求
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(resp, req)

			// 验证状态码
			assert.Equal(t, tt.expectedStatus, resp.Code)

			// 验证响应体
			var actualBody map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualBody)
		})
	}
}

func TestErrorHandler_LoggingBehavior(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试日志记录器
	core, logs := observer.New(zap.InfoLevel)
	testLogger := logger.New(zap.New(core))

	// 创建路由器
	router := gin.New()
	router.Use(ErrorHandler(testLogger))
	router.GET("/test", func(c *gin.Context) {
		c.Error(&gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  errors.New("test error"),
		})
	})

	// 创建测试请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(resp, req)

	// 验证日志记录
	assert.Equal(t, 1, logs.Len())
	logEntry := logs.All()[0]
	assert.Equal(t, "内部错误", logEntry.Message)
	assert.Equal(t, "/test", logEntry.Context[0].String)
	assert.Equal(t, "GET", logEntry.Context[1].String)
	assert.Equal(t, "test error", logEntry.Context[2].Interface.(error).Error())
}
