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

type CredentialController struct {
	credService services.CredentialInterface
}

func NewCredentialController(credService services.CredentialInterface) *CredentialController {
	return &CredentialController{
		credService: credService,
	}
}

func (cc *CredentialController) Write(c *gin.Context) {
	var payload models.UserPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		res := utilities.RegistrationResponse{
			Message: "validation error",
			Err:     err,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*1)
	defer cancel()

	id, err := cc.credService.Write(ctx, payload)
	if err != nil {
		res := utilities.RegistrationResponse{
			Message: "failed to create user",
			Err:     err,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utilities.RegistrationResponse{
		ID:      id,
		Message: "success",
	}

	c.JSON(http.StatusCreated, res)
}
