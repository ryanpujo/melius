package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/melius/internal/adapter"
	"github.com/ryanpujo/melius/internal/jwttoken"
)
var router = gin.Default()
var apiV1 = router.Group("/api/v1")
var protectedApiV1 = apiV1.Group("/auth")
// SetupRoutes initializes and returns a Gin engine with defined routes.
func SetupRoutes(handlers *adapter.Adapter, auth jwttoken.Authenticator) *gin.Engine {
	protectedApiV1.Use(auth.JWTAuthMiddleware())
	// Define a simple GET route
	protectedApiV1.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, World!")
	})

	addressRoute(handlers.AddressController)
	router.POST("/regis", handlers.CredentialController.Write)
	router.POST("/login", handlers.CredentialController.Login)

	return router
}
