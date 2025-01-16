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

type AddressController interface {
	SaveCountry(c *gin.Context)
	SaveState(c *gin.Context)
}

type addressController struct {
	addressService services.AddressService
}

func NewAddressController(addressService services.AddressService) *addressController {
	return &addressController{
		addressService: addressService,
	}
}

func (ac *addressController) SaveCountry(c *gin.Context) {
	var json models.Country
	if err := c.ShouldBindJSON(&json); err != nil {
		res := utilities.Response{
			Message: "Validation Error",
			Err:     err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*1)
	defer cancel()

	id, err := ac.addressService.SaveCountry(ctx, json)
	if err != nil {
		res := utilities.Response{
			Message: "failed to record the country",
			Err:     err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utilities.Response{
		Message: "Success",
		ID:      id,
	}
	c.JSON(http.StatusCreated, res)
}

type uriBind struct {
	ID uint `uri:"id"  binding:"required"`
}

func (ac *addressController) SaveState(c *gin.Context) {
	var json models.State
	var uri uriBind

	if err := c.ShouldBindUri(&uri); err != nil {
		res := utilities.Response{
			Message: "there is no country to be associated with this state",
			Err:     err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		res := utilities.Response{
			Message: "Validation Error",
			Err:     err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*1)
	defer cancel()

	id, err := ac.addressService.SaveState(ctx, json, uri.ID)
	if err != nil {
		res := utilities.Response{
			Message: "failed to record the state",
			Err:     err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utilities.Response{
		Message: "Success",
		ID:      id,
	}
	c.JSON(http.StatusCreated, res)
}
