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

// AddressController defines the interface for address-related operations.
type AddressController interface {
	SaveCountry(c *gin.Context)
	SaveState(c *gin.Context)
}

// addressController implements the AddressController interface.
type addressController struct {
	addressService services.AddressService
}

// NewAddressController creates a new instance of addressController.
func NewAddressController(addressService services.AddressService) *addressController {
	return &addressController{
		addressService: addressService,
	}
}

// SaveCountry handles the HTTP request to save a new country.
// @param c *gin.Context - The Gin context containing request data.
func (ac *addressController) SaveCountry(c *gin.Context) {
	var country models.Country

	// Validate the JSON request body.
	if err := c.ShouldBindJSON(&country); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Validation Error", utilities.WithErr(err.Error())),
		)
		return
	}

	// Set a context with a timeout.
	ctx, cancel := context.WithTimeout(c, time.Second*1)
	defer cancel()

	// Call the service layer to save the country.
	id, err := ac.addressService.SaveCountry(ctx, country)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Failed to record the country", utilities.WithErr(err.Error())),
		)
		return
	}

	// Respond with success.
	c.JSON(
		http.StatusCreated,
		utilities.NewResponse("Success", utilities.WithID(id)),
	)
}

// SaveState handles the HTTP request to save a new state associated with a country.
// @param c *gin.Context - The Gin context containing request data.
func (ac *addressController) SaveState(c *gin.Context) {
	var state models.State
	var uri uriBind

	// Validate the URI parameters.
	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("No country associated with this state", utilities.WithErr(err.Error())),
		)
		return
	}

	// Validate the JSON request body.
	if err := c.ShouldBindJSON(&state); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Validation Error", utilities.WithErr(err.Error())),
		)
		return
	}

	// Set a context with a timeout.
	ctx, cancel := context.WithTimeout(c, time.Second*1)
	defer cancel()

	// Call the service layer to save the state.
	id, err := ac.addressService.SaveState(ctx, state, uri.ID)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Failed to record the state", utilities.WithErr(err.Error())),
		)
		return
	}

	// Respond with success.
	c.JSON(
		http.StatusCreated,
		utilities.NewResponse("Success", utilities.WithID(id)),
	)
}

// uriBind represents the structure for URI parameters.
type uriBind struct {
	ID uint `uri:"id" binding:"required"`
}
