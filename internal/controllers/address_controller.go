package controllers

import (
	"context"
	"database/sql"
	"errors"
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
	SaveCity(c *gin.Context)
	SaveAddress(c *gin.Context)

	GetCountries(c *gin.Context)
	GetCityByID(c *gin.Context)
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
	var country models.CountryPayload

	// Validate the JSON request body.
	if err := c.ShouldBindJSON(&country); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Validation Error", utilities.WithErr(err.Error())),
		)
		return
	}

	// Set a context with a timeout.
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*1)
	defer cancel()

	// Call the service layer to save the country.
	createdCountry, err := ac.addressService.SaveCountry(ctx, &country)
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
		utilities.NewResponse("Success", utilities.WithCountry(createdCountry)),
	)
}

// SaveState handles the HTTP request to save a new state associated with a country.
// @param c *gin.Context - The Gin context containing request data.
func (ac *addressController) SaveState(c *gin.Context) {
	var state models.StatePayload
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*1)
	defer cancel()

	// Call the service layer to save the state.
	createdState, err := ac.addressService.SaveState(ctx, &state, uri.ID)
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
		utilities.NewResponse("Success", utilities.WithState(createdState)),
	)
}

func (ac *addressController) SaveCity(c *gin.Context) {
	var city models.CityPayload
	var uri uriBind

	// Validate the URI parameters.
	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("No state associated with this city", utilities.WithErr(err.Error())),
		)
		return
	}

	// Validate the JSON request body.
	if err := c.ShouldBindJSON(&city); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Validation Error", utilities.WithErr(err.Error())),
		)
		return
	}

	// Set a context with a timeout.
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*1)
	defer cancel()

	// Call the service layer to save the city.
	createdCity, err := ac.addressService.SaveCity(ctx, &city, uri.ID)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Failed to record the city", utilities.WithErr(err.Error())),
		)
		return
	}

	// Respond with success.
	c.JSON(
		http.StatusCreated,
		utilities.NewResponse("Success", utilities.WithCity(createdCity)),
	)
}

func (ac *addressController) SaveAddress(c *gin.Context) {
	var address models.AddressPayload

	// Validate the JSON request body.
	if err := c.ShouldBindJSON(&address); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Validation Error", utilities.WithErr(err.Error())),
		)
		return
	}

	// Set a context with a timeout.
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*1)
	defer cancel()

	// Call the service layer to save the address.
	createdAddress, err := ac.addressService.SaveAddress(ctx, &address)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Failed to record the address", utilities.WithErr(err.Error())),
		)
		return
	}

	// Respond with success.
	c.JSON(
		http.StatusCreated,
		utilities.NewResponse("Success", utilities.WithAddress(createdAddress)),
	)
}

func (ac *addressController) GetCountries(c *gin.Context) {
	// Set a context with a timeout.
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*1)
	defer cancel()

	countries, err := ac.addressService.GetCountries(ctx)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utilities.NewResponse("Failed to get countries", utilities.WithErr(err.Error())),
		)
		return
	}

	c.JSON(http.StatusOK, utilities.NewResponse("success", utilities.WithCountries(countries)))
}

func (ac *addressController) GetCityByID(c *gin.Context) {
	var uri uriBind

	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			utilities.NewResponse(
				"There was a problem with your request. Please double-check your input and try again.",
				utilities.WithErr(err.Error()),
			),
		)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*1)
	defer cancel()

	city, err := ac.addressService.GetCityByID(ctx, uri.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				utilities.NewResponse(
					"We're sorry, but we couldn't find a city with that information. Please check your input and try again.",
					utilities.WithErr(err.Error()),
				),
			)
			return
		}
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utilities.NewResponse(
				"An unexpected error occurred. Please try again later.",
				utilities.WithErr(err.Error()),
			),
		)
		return
	}

	c.JSON(http.StatusOK, utilities.NewResponse("success", utilities.WithCity(city)))
}

// uriBind represents the structure for URI parameters.
type uriBind struct {
	ID uint `uri:"id" binding:"required"`
}
