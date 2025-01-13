package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/melius/internal/adapter"
	"github.com/ryanpujo/melius/internal/jwttoken"
)
var router = gin.Default()
var protected = router.Group("/auth")
// SetupRoutes initializes and returns a Gin engine with defined routes.
func SetupRoutes(handlers *adapter.Adapter, auth jwttoken.Authenticator) *gin.Engine {
	protected.Use(auth.JWTAuthMiddleware())
	// Define a simple GET route
	protected.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, World!")
	})

	addressRoute(handlers.AddressController)
	router.POST("/regis", handlers.CredentialController.Write)
	router.POST("/login", handlers.CredentialController.Login)

	return router
}
