package services

import (
	"context"
	"errors"

	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/repositories"
)

type AddressService interface {
	SaveCountry(ctx context.Context, country *models.CountryPayload) (*models.Country, error)
	SaveState(ctx context.Context, state *models.StatePayload, countryID uint) (*models.State, error)
	SaveCity(ctx context.Context, city *models.CityPayload, stateID uint) (*models.City, error)
	SaveAddress(ctx context.Context, address *models.AddressPayload) (*models.Address, error)

	GetCountries(ctx context.Context) ([]*models.Country, error)
	GetCityByID(ctx context.Context, cityID uint) (*models.City, error)
}

type addressService struct {
	addressRepo repositories.AddressRepo
}

func NewAddressService(addressRepo repositories.AddressRepo) *addressService {
	return &addressService{
		addressRepo: addressRepo,
	}
}

// SaveCountry saves a country to the repository.
// It takes the context and a Country model as parameters.
// Returns the generated country ID and any error encountered during the operation.
func (as *addressService) SaveCountry(ctx context.Context, country *models.CountryPayload) (*models.Country, error) {
	return as.addressRepo.SaveCountry(ctx, country, nil)
}

func (as *addressService) GetCountries(ctx context.Context) ([]*models.Country, error) {
	return as.addressRepo.GetCountries(ctx)
}

// SaveState saves a state to the repository.
// It takes the context, a State model, and a countryID as parameters.
// Validates that the countryID is not empty before proceeding.
// Returns the generated state ID and any error encountered during the operation.
func (as *addressService) SaveState(ctx context.Context, state *models.StatePayload, countryID uint) (*models.State, error) {
	if countryID == 0 {
		return nil, errors.New("country ID cannot be empty")
	}
	return as.addressRepo.SaveState(ctx, state, countryID, nil)
}

// SaveCity saves a city to the repository.
// It takes the context, a City model, and a stateID as parameters.
// Validates that the stateID is not empty before proceeding.
// Returns the generated city ID and any error encountered during the operation.
func (as *addressService) SaveCity(ctx context.Context, city *models.CityPayload, stateID uint) (*models.City, error) {
	if stateID == 0 {
		return nil, errors.New("state ID cannot be empty")
	}
	return as.addressRepo.SaveCity(ctx, city, stateID, nil)
}

func (as *addressService) GetCityByID(ctx context.Context, cityID uint) (*models.City, error) {
	if cityID == 0 {
		return nil, errors.New("city id cannot be empty")
	}

	return as.addressRepo.GetCityByID(ctx, cityID)
}

// SaveAddress saves an address to the repository.
// It takes the context, an Address model, and a cityID as parameters.
// Validates that the cityID is not empty before proceeding.
// Returns the generated address ID and any error encountered during the operation.
func (as *addressService) SaveAddress(ctx context.Context, address *models.AddressPayload) (*models.Address, error) {
	if address.CityID == 0 {
		return nil, errors.New("city ID cannot be empty")
	}
	return as.addressRepo.SaveAddress(ctx, address, nil)
}