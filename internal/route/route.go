package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/melius/internal/adapter"
)

// SetupRoutes initializes and returns a Gin engine with defined routes.
func SetupRoutes(handlers *adapter.Adapter) *gin.Engine {
	router := gin.Default()

	// Define a simple GET route
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, World!")
	})

	router.POST("/regis", handlers.CredentialController.Write)

	return router
}
