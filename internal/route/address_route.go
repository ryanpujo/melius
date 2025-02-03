package route

import (
	"github.com/ryanpujo/melius/internal/controllers"
)

func addressRoute(addressController controllers.AddressController) {
	protectedApiV1.POST("/country", addressController.SaveCountry)
	protectedApiV1.GET("/country", addressController.GetCountries)

	protectedApiV1.POST("/state/:id", addressController.SaveState)

	protectedApiV1.POST("/city/:id", addressController.SaveCity)
	protectedApiV1.GET("/city/:id", addressController.GetCityByID)

	protectedApiV1.POST("/address", addressController.SaveAddress)
}
