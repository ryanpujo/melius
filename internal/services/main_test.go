package services_test

import (
	"os"
	"testing"

	"github.com/ryanpujo/melius/internal/services"
)

var (
	credService    services.CredentialService
	addressService services.AddressService
	crm            *CredRepoMock
	arm            *addressRepoMock
)

func TestMain(m *testing.M) {
	crm = new(CredRepoMock)
	arm = new(addressRepoMock)
	credService = *services.NewCredentialService(crm)
	addressService = services.NewAddressService(arm)
	os.Exit(m.Run())
}
