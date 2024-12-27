package middleware

import (
	"net/http"
	"strings"

	"blockchain-shop/pkg/errors"
	"blockchain-shop/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	secretKey string
	issuer    string
}

func NewJWTMiddleware(secretKey string, issuer string) (*JWTMiddleware, error) {
	return &JWTMiddleware{secretKey: secretKey, issuer: issuer}, nil
}

func (m *JWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Error(gin.Error{
				Type: gin.ErrorTypePublic,
				Err:  errors.ErrUnauthorized,
				Meta: http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Error(gin.Error{
				Type: gin.ErrorTypePublic,
				Err:  errors.ErrUnauthorized,
				Meta: http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwt.ParseToken(token, m.secretKey)
		if err != nil {
			c.Error(gin.Error{
				Type: gin.ErrorTypePublic,
				Err:  errors.ErrUnauthorized,
				Meta: http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
