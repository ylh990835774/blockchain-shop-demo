package service

import (
	"time"

	"github.com/ylh990835774/blockchain-shop-demo/pkg/jwt"
)

type JWTService struct {
	secretKey string
	issuer    string
}

func NewJWTService(secretKey string, issuer string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (s *JWTService) GenerateToken(userID int64) (string, error) {
	return jwt.GenerateToken(userID, s.secretKey, s.issuer, 24*time.Hour)
}
