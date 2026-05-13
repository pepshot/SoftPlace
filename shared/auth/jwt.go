package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pepshot/SoftPlace/shared/config"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret   []byte
	lifetime time.Duration
}

func NewJWTManager(cfg config.AuthConfig) *JWTManager {
	lifetime := time.Duration(cfg.JWTLifetimeHours) * time.Hour
	if cfg.JWTLifetimeHours <= 0 {
		lifetime = 24 * time.Hour
	}

	return &JWTManager{
		secret:   []byte(cfg.JWTSecret),
		lifetime: lifetime,
	}
}

func (m *JWTManager) GenerateToken(userID string, role string) (string, error) {
	now := time.Now()
	claims := &Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.lifetime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
