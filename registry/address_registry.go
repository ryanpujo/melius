package registry

import (
	"github.com/ryanpujo/melius/internal/controllers"
	"github.com/ryanpujo/melius/internal/repositories"
	"github.com/ryanpujo/melius/internal/services"
)

func (r *Registry) NewAddressRepository() repositories.AddressRepo {
	return repositories.NewAddressRepo(r.db)
}

func (r *Registry) NewAddressService() services.AddressService {
	return services.NewAddressService(r.NewAddressRepository())
}

func (r *Registry) NewAddressController() controllers.AddressController {
	return controllers.NewAddressController(r.NewAddressService())
}
