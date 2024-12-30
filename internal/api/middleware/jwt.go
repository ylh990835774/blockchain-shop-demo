package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ylh990835774/blockchain-shop-demo/internal/service"
)

// JWTMiddleware 是JWT中间件的结构体
type JWTMiddleware struct {
	jwtService service.IJWTService
}

// NewJWTMiddleware 创建一个新的JWT中间件
func NewJWTMiddleware(jwtService service.IJWTService) *JWTMiddleware {
	if jwtService == nil {
		panic("jwt service cannot be nil")
	}

	return &JWTMiddleware{
		jwtService: jwtService,
	}
}

// MiddlewareFunc 返回JWT中间件的处理函数
func (m *JWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "未授权的访问",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "无效的认证格式",
			})
			return
		}

		token := parts[1]
		userID, err := m.jwtService.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "无效的认证令牌" + err.Error() + "(token:" + token + ")",
			})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
