package jwttoken

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/melius/config"
	"github.com/ryanpujo/melius/internal/utilities"
)

type TokenVerifier interface {
	VerifyToken(token string) (*jwt.Token, error)
}

type tokenVerif struct {
	jwtKey string
}

func (tv *tokenVerif) VerifyToken(tokenString string) (*jwt.Token, error) {
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

func NewTokenVerif() *tokenVerif {
	return &tokenVerif{
		jwtKey: config.Config().JWTKey,
	}
}

type Authenticator interface {
	GenerateJWT(username string) (string, error)
	JWTAuthMiddleware() gin.HandlerFunc
}

type JWTAuth struct {
	tokenVerifier TokenVerifier
}

var jwtAuth *JWTAuth

func GetJWTAuth(tokenVerifier TokenVerifier) *JWTAuth {
	if jwtAuth == nil {
		jwtAuth = &JWTAuth{
			tokenVerifier: tokenVerifier,
		}
	}
	return jwtAuth
}

func (auth *JWTAuth) GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * time.Duration(config.Config().JWTConfig.EXP)).Unix(),
		"aud":      config.Config().JWTConfig.AUD,
		"iss":      config.Config().JWTConfig.ISS,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Config().JWTKey))
}

func (auth *JWTAuth) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}
		tokenString, _ := strings.CutPrefix(authHeader, "Bearer")

		token, err := auth.tokenVerifier.VerifyToken(tokenString)
		if err != nil {
			res := utilities.Response{
				Message: "authentication failed",
				Err: err.Error(),
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		c.Set("username", claims["username"])
		c.Next()
	}
}
