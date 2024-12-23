package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/services"
	"github.com/ryanpujo/melius/internal/utilities"
)

// CredentialController handles user authentication and registration.
type CredentialController struct {
	credService services.CredentialInterface
}

// NewCredentialController initializes a new CredentialController with the provided credential service.
func NewCredentialController(credService services.CredentialInterface) *CredentialController {
	return &CredentialController{
		credService: credService,
	}
}

// Write handles user registration.
// It validates the payload, calls the credential service to create a user,
// and returns the created user ID or an error response.
func (cc *CredentialController) Write(c *gin.Context) {
	var payload models.UserPayload

	// Bind and validate JSON payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utilities.RegistrationResponse{
			Message: "Validation error",
			Err:     err.Error(),
		})
		return
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(c, time.Second*1)
	defer cancel()

	// Call service to create user
	id, err := cc.credService.Write(ctx, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utilities.RegistrationResponse{
			Message: "Failed to create user",
			Err:     err.Error(),
		})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, utilities.RegistrationResponse{
		ID:      id,
		Message: "User created successfully",
	})
}

// Login handles user authentication.
// It validates the payload, calls the credential service to authenticate the user,
// and returns a JWT token or an error response.
func (cc *CredentialController) Login(c *gin.Context) {
	var payload models.LoginPayload

	// Bind and validate JSON payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utilities.RegistrationResponse{
			Message: "Validation error",
			Err:     err.Error(),
		})
		return
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(c, time.Second*1)
	defer cancel()

	// Call service to authenticate user
	jwt, err := cc.credService.Login(ctx, &payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utilities.RegistrationResponse{
			Message: "Login failed",
			Err:     err.Error(),
		})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, utilities.RegistrationResponse{
		Message: "Login successful",
		Token:   jwt,
	})
}
