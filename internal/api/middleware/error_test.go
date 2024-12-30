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
)

func TestErrorHandler(t *testing.T) {
	// 初始化日志配置
	err := logger.Setup(&logger.Config{
		Level:    "info",
		Format:   "console",
		Console:  true,
		Filename: "", // 测试时不写入文件
	})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupRouter    func(*gin.Engine)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "public error with status code",
			setupRouter: func(r *gin.Engine) {
				r.Use(ErrorHandler())
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
				"data":    nil,
			},
		},
		{
			name: "not found error",
			setupRouter: func(r *gin.Engine) {
				r.Use(ErrorHandler())
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
				"data":    nil,
			},
		},
		{
			name: "unauthorized error",
			setupRouter: func(r *gin.Engine) {
				r.Use(ErrorHandler())
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
				"data":    nil,
			},
		},
		{
			name: "invalid input error",
			setupRouter: func(r *gin.Engine) {
				r.Use(ErrorHandler())
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
				"data":    nil,
			},
		},
		{
			name: "duplicate entry error",
			setupRouter: func(r *gin.Engine) {
				r.Use(ErrorHandler())
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
				"data":    nil,
			},
		},
		{
			name: "unknown error",
			setupRouter: func(r *gin.Engine) {
				r.Use(ErrorHandler())
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
				"data":    nil,
			},
		},
		{
			name: "no error",
			setupRouter: func(r *gin.Engine) {
				r.Use(ErrorHandler())
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
			// 创建路由器
			router := gin.New()
			tt.setupRouter(router)

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
