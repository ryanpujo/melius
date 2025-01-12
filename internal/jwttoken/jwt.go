package jwttoken

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/melius/config"
)

type JWTAuth struct {
	EXP int64
	AUD any
	ISS string
}

var jwtAuth *JWTAuth

func GetJWTAuth() *JWTAuth {
	if jwtAuth == nil {
		jwtAuth = &JWTAuth{
			EXP: int64(config.Config().JWTConfig.EXP), // Short expiration time
			AUD: config.Config().JWTConfig.AUD,
			ISS: config.Config().JWTConfig.ISS,
		}
	}
	return jwtAuth
}

func (auth *JWTAuth) GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      auth.EXP,
		"aud":      auth.AUD,
		"iss":      auth.ISS,
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

		token, err := jwt.Parse(strings.TrimSpace(tokenString), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Config().JWTKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Next()
	}
}
