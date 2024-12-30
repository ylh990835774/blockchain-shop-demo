package service

import (
	"time"

	"github.com/ylh990835774/blockchain-shop-demo/pkg/jwt"
)

type jwtService struct {
	secretKey      string
	issuer         string
	expireDuration time.Duration
}

// 确保jwtService实现了IJWTService接口
var _ IJWTService = (*jwtService)(nil)

// NewJWTService 创建一个新的JWT服务
func NewJWTService(secretKey, issuer string, expireDuration time.Duration) IJWTService {
	return &jwtService{
		secretKey:      secretKey,
		issuer:         issuer,
		expireDuration: expireDuration,
	}
}

// GenerateToken 生成JWT令牌
func (s *jwtService) GenerateToken(userID int64) (string, error) {
	return jwt.GenerateToken(userID, s.secretKey, s.issuer, s.expireDuration)
}

// ParseToken 解析JWT令牌
func (s *jwtService) ParseToken(token string) (int64, error) {
	claims, err := jwt.ParseToken(token, s.secretKey)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
