package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 包含自定义的用户ID和标准JWT声明
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// ParseToken 解析并验证JWT令牌
func ParseToken(tokenString, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// GenerateToken 生成新的JWT令牌
func GenerateToken(userID int64, secretKey, issuer string, expireDuration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
