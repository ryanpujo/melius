package jwttoken

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/melius/config"
	"github.com/ryanpujo/melius/internal/utilities"
)

// Authenticator defines an interface for generating JWT tokens and handling authentication middleware.
type Authenticator interface {
	GenerateJWT(username string) (string, error)
	JWTAuthMiddleware() gin.HandlerFunc
}

// JWTAuth is an implementation of the Authenticator interface.
type JWTAuth struct {
	tokenVerifier TokenVerifier
}

var jwtAuthInstance *JWTAuth

// GetJWTAuth returns a singleton instance of JWTAuth.
// It initializes the instance if it does not already exist.
func GetJWTAuth(tokenVerifier TokenVerifier) *JWTAuth {
	if jwtAuthInstance == nil {
		jwtAuthInstance = &JWTAuth{
			tokenVerifier: tokenVerifier,
		}
	}
	return jwtAuthInstance
}

// GenerateJWT generates a signed JWT token with the provided username and configured claims.
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

// JWTAuthMiddleware is a Gin middleware for authenticating requests using JWT tokens.
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
				Err:     err.Error(),
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("username", claims["username"])
		c.Next()
	}
}
