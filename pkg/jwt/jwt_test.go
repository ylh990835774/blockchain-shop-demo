package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		secretKey   string
		issuer      string
		duration    time.Duration
		expectErr   bool
		modifyToken func(string) string
	}{
		{
			name:      "valid token",
			userID:    123,
			secretKey: "test-secret-key",
			issuer:    "test-issuer",
			duration:  time.Hour,
			expectErr: false,
		},
		{
			name:      "expired token",
			userID:    123,
			secretKey: "test-secret-key",
			issuer:    "test-issuer",
			duration:  -time.Hour, // 过期的token
			expectErr: true,
		},
		{
			name:      "invalid signature",
			userID:    123,
			secretKey: "test-secret-key",
			issuer:    "test-issuer",
			duration:  time.Hour,
			expectErr: true,
			modifyToken: func(token string) string {
				return token + "invalid"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 生���token
			token, err := GenerateToken(tt.userID, tt.secretKey, tt.issuer, tt.duration)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// 如果需要修改token
			if tt.modifyToken != nil {
				token = tt.modifyToken(token)
			}

			// 解析token
			claims, err := ParseToken(token, tt.secretKey)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, claims)
			assert.Equal(t, tt.userID, claims.UserID)
			assert.Equal(t, tt.issuer, claims.Issuer)
		})
	}
}

func TestGenerateToken_InvalidParams(t *testing.T) {
	// 测试空密钥
	token, err := GenerateToken(123, "", "test", time.Hour)
	assert.Error(t, err, "should return error with empty secret key")
	assert.Empty(t, token, "should return empty token with empty secret key")
}

func TestParseToken_InvalidParams(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		secretKey string
	}{
		{
			name:      "empty token",
			token:     "",
			secretKey: "test-key",
		},
		{
			name:      "invalid token format",
			token:     "invalid.token.format",
			secretKey: "test-key",
		},
		{
			name:      "empty secret key",
			token:     "valid.token.format",
			secretKey: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token, tt.secretKey)
			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}
