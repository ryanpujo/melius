package jwttoken

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/melius/config"
)

// TokenVerifier defines an interface for verifying JWT tokens.
type TokenVerifier interface {
	VerifyToken(token string) (*jwt.Token, error)
}

// tokenVerifier is an implementation of the TokenVerifier interface.
type tokenVerifier struct {
	jwtKey string
}

// VerifyToken verifies the validity of a JWT token string.
// It checks the token's signature, validity, and expiration.
func (tv *tokenVerifier) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(strings.TrimSpace(tokenString), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(tv.jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
		return nil, errors.New("token is expired")
	}
	return token, nil
}

// NewTokenVerifier creates a new tokenVerifier instance using the configured JWT key.
func NewTokenVerifier() *tokenVerifier {
	return &tokenVerifier{
		jwtKey: config.Config().JWTKey,
	}
}
