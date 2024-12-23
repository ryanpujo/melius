package registry

import (
	"github.com/ryanpujo/melius/internal/controllers"
	"github.com/ryanpujo/melius/internal/repositories"
	"github.com/ryanpujo/melius/internal/services"
)

func (r *Registry) GetCredentialRepo() repositories.CredentialInterface {
	return repositories.NewCredentialRepo(r.db)
}

func (r *Registry) GetCredentialService() services.CredentialInterface {
	return services.NewCredentialService(r.GetCredentialRepo())
}

func (r *Registry) GetCredentialController() *controllers.CredentialController {
	return controllers.NewCredentialController(r.GetCredentialService())
}
