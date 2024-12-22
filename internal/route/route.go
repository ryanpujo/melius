package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/melius/internal/adapter"
	"github.com/ryanpujo/melius/internal/jwttoken"
)

// SetupRoutes initializes and returns a Gin engine with defined routes.
func SetupRoutes(handlers *adapter.Adapter) *gin.Engine {
	router := gin.Default()
	protected := router.Group("/auth")
	protected.Use(jwttoken.JWTAuthMiddleware())
	// Define a simple GET route
	protected.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, World!")
	})

	router.POST("/regis", handlers.CredentialController.Write)
	router.POST("/login", handlers.CredentialController.Login)

	return router
}
