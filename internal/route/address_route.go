package route

import (
	"github.com/ryanpujo/melius/internal/controllers"
)

func addressRoute(addressController controllers.AddressController) {
	protected.POST("/country", addressController.SaveCountry)
}
