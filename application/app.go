package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ryanpujo/melius/config"
)

// Application represents the server configuration.
type Application struct {
	Port    int          // Port on which the server will run
	Handler http.Handler // HTTP handler (e.g., routes)
}

// NewApp initializes a new Application with the given handler.
func NewApp(handler http.Handler) *Application {
	conf := config.Config()
	return &Application{
		Port:    conf.Port,
		Handler: handler,
	}
}

// Serve starts the HTTP server with defined timeouts.
func (app *Application) Serve() error {
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", app.Port),
		Handler:           app.Handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	fmt.Printf("Server is running on port %d\n", app.Port)
	return server.ListenAndServe()
}
