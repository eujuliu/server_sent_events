package utils

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		slog.Debug(fmt.Sprintf("JWT validation failed: %s", err))
		return nil, fmt.Errorf("token validation failed")
	}

	claims, ok := token.Claims.(*Claims)

	if !ok && !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if expired := isTokenExpired(claims); expired {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func isTokenExpired(claims *Claims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}
