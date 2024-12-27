package service

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	now := time.Now()
	claims := jwt.StandardClaims{
		ExpiresAt: now.Add(24 * time.Hour).Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    s.issuer,
		Subject:   strconv.FormatInt(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}
