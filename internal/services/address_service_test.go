package services_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/utilities"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type addressRepoMock struct {
	mock.Mock
}

func (ar *addressRepoMock) SaveCountry(ctx context.Context, country *models.CountryPayload, tx *sql.Tx) (*models.Country, error) {
	args := ar.Called(ctx, country, tx)
	return args.Get(0).(*models.Country), args.Error(1)
}

func (ar *addressRepoMock) SaveState(ctx context.Context, state *models.StatePayload, countryID uint, tx *sql.Tx) (*models.State, error) {
	args := ar.Called(ctx, state, countryID, tx)
	return args.Get(0).(*models.State), args.Error(1)
}

func (ar *addressRepoMock) SaveCity(ctx context.Context, city *models.CityPayload, stateID uint, tx *sql.Tx) (*models.City, error) {
	args := ar.Called(ctx, city, stateID, tx)
	return args.Get(0).(*models.City), args.Error(1)
}

func (ar *addressRepoMock) SaveAddress(ctx context.Context, address *models.AddressPayload, tx *sql.Tx) (*models.Address, error) {
	args := ar.Called(ctx, address, tx)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (ar *addressRepoMock) GetCountries(ctx context.Context) ([]*models.Country, error) {
	args := ar.Called(ctx)
	return args.Get(0).([]*models.Country), args.Error(1)
}

func (ar *addressRepoMock) GetCityByID(ctx context.Context, cityID uint) (*models.City, error) {
	args := ar.Called(ctx, cityID)
	return args.Get(0).(*models.City), args.Error(1)
}

func (ar *addressRepoMock) CreateUserAddress(ctx context.Context, address *models.AddressPayload, userID uint) (*models.Address, error) {
	args := ar.Called(ctx, address, userID)
	return args.Get(0).(*models.Address), args.Error(1)
}

var (
	countryPayload = &models.CountryPayload{
		Name:  "Indonesia",
	}
	country = &models.Country{
		ID:   1,
		Name: "Indonesia",
	}

	statePayload = models.StatePayload{
		Name: "Jakarta",
	}
	state = &models.State{
		ID:   1,
		Name: "Jakarta",
	}

	cityPayload = models.CityPayload{
		Name:    "Jakarta Timur",
	}
	city = &models.City{
		ID:   1,
		Name: "Jakarta Timur",
	}

	addressPayload = models.AddressPayload{
		AddressLine: "jl. mayjen sutoyo kel. cawang kec kramat jati rt.007/011",
		PostalCode:  "12630",
		IsMain:      true,
		CityID:      uint(1),
	}
	address = &models.Address{
		AddressLine: "jl. mayjen sutoyo kel. cawang kec kramat jati rt.007/011",
		PostalCode:  "12630",
		IsMain:      true,
		City:        city,
	}
)

func TestSaveCountry(t *testing.T) {
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualCountry *models.Country, err error)
	}{
		"succes": {
			arrange: func() {
				arm.On("SaveCountry", mock.Anything, countryPayload, (*sql.Tx)(nil)).Return(country, nil).Once()
			},
			assert: func(t *testing.T, actualCountry *models.Country, err error) {
				require.NoError(t, err)
				require.Equal(t, country, actualCountry)
				arm.AssertCalled(t, "SaveCountry", context.Background(), countryPayload, (*sql.Tx)(nil))
			},
		},
		"failed": {
			arrange: func() {
				arm.On("SaveCountry", mock.Anything, countryPayload, (*sql.Tx)(nil)).Return((*models.Country)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualCountry *models.Country, err error) {
				require.Error(t, err)
				require.Zero(t, actualCountry)
				arm.AssertCalled(t, "SaveCountry", context.Background(), countryPayload, (*sql.Tx)(nil))
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			country, err := addressService.SaveCountry(context.Background(), countryPayload)

			v.assert(t, country, err)
		})
	}
}

func TestSaveState(t *testing.T) {
	tableTest := map[string]struct {
		countryID uint
		arrange   func()
		assert    func(t *testing.T, actualState *models.State, err error)
	}{
		"success": {
			countryID: 1,
			arrange: func() {
				arm.On("SaveState", mock.Anything, &statePayload, uint(1), (*sql.Tx)(nil)).Return(state, nil).Once()
			},
			assert: func(t *testing.T, actualState *models.State, err error) {
				require.NoError(t, err)
				require.Equal(t, state, actualState)
				arm.AssertCalled(t, "SaveState", context.Background(), &statePayload, uint(1), (*sql.Tx)(nil))
			},
		},
		"failed": {
			countryID: 1,
			arrange: func() {
				arm.On("SaveState", mock.Anything, &statePayload, uint(1), (*sql.Tx)(nil)).
					Return((*models.State)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualState *models.State, err error) {
				require.Error(t, err)
				require.Zero(t, actualState)
				arm.AssertCalled(t, "SaveState", context.Background(), &statePayload, uint(1), (*sql.Tx)(nil))
			},
		},
		"empty country id": {
			countryID: 0,
			arrange:   func() {},
			assert: func(t *testing.T, actualState *models.State, err error) {
				var appErr *utilities.AppError
				require.Error(t, err)
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, utilities.ValidationError, appErr.Code)
				require.Equal(t, "country ID can't be empty", appErr.Message)
				require.Zero(t, actualState)
				arm.AssertNotCalled(t, "SaveState")
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			state, err := addressService.SaveState(context.Background(), &statePayload, v.countryID)

			v.assert(t, state, err)
		})
	}
}

func TestSaveCity(t *testing.T) {
	tableTest := map[string]struct {
		stateID uint
		arrange func()
		assert  func(t *testing.T, actualCity *models.City, err error)
	}{
		"success": {
			stateID: 1,
			arrange: func() {
				arm.On("SaveCity", mock.Anything, &cityPayload, uint(1), (*sql.Tx)(nil)).Return(city, nil).Once()
			},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				require.NoError(t, err)
				require.Equal(t, city, actualCity)
				arm.AssertCalled(t, "SaveCity", context.Background(), &cityPayload, uint(1), (*sql.Tx)(nil))
			},
		},
		"failed": {
			stateID: 1,
			arrange: func() {
				arm.On("SaveCity", mock.Anything, &cityPayload, uint(1), (*sql.Tx)(nil)).
					Return((*models.City)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				require.Error(t, err)
				require.Zero(t, actualCity)
				arm.AssertCalled(t, "SaveCity", context.Background(), &cityPayload, uint(1), (*sql.Tx)(nil))
			},
		},
		"empty country id": {
			stateID: 0,
			arrange: func() {},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				var appErr *utilities.AppError
				require.Error(t, err)
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, utilities.ValidationError, appErr.Code)
				require.Equal(t, "state ID can't be empty", appErr.Message)
				require.Zero(t, actualCity)
				arm.AssertNotCalled(t, "SaveCity")
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			city, err := addressService.SaveCity(context.Background(), &cityPayload, v.stateID)

			v.assert(t, city, err)
		})
	}
}

func TestSaveAddress(t *testing.T) {
	tableTest := map[string]struct {
		cityID  uint
		arrange func()
		assert  func(t *testing.T, actualAddress *models.Address, err error)
	}{
		"success": {
			cityID: 1,
			arrange: func() {
				arm.On("SaveAddress", mock.Anything, &addressPayload, (*sql.Tx)(nil)).
					Return(address, nil).Once()
			},
			assert: func(t *testing.T, actualAddress *models.Address, err error) {
				require.NoError(t, err)
				require.Equal(t, address, actualAddress)
				arm.AssertCalled(t, "SaveAddress", context.Background(), &addressPayload, (*sql.Tx)(nil))
			},
		},
		"failed": {
			cityID: 1,
			arrange: func() {
				arm.On("SaveAddress", mock.Anything, &addressPayload, (*sql.Tx)(nil)).
					Return((*models.Address)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualAddress *models.Address, err error) {
				require.Error(t, err)
				require.Zero(t, actualAddress)
				arm.AssertCalled(t, "SaveAddress", context.Background(), &addressPayload, (*sql.Tx)(nil))
			},
		},
		"empty country id": {
			cityID:  0,
			arrange: func() {},
			assert: func(t *testing.T, actualAddress *models.Address, err error) {
				var appErr *utilities.AppError
				require.Error(t, err)
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, utilities.ValidationError, appErr.Code)
				require.Equal(t, "city ID can't be empty", appErr.Message)
				require.Zero(t, actualAddress)
				arm.AssertNotCalled(t, "SaveAddress")
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			addressPayload.CityID = v.cityID
			address, err := addressService.SaveAddress(context.Background(), &addressPayload)

			v.assert(t, address, err)
		})
	}
}

func TestGetCountries(t *testing.T) {
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualCountries []*models.Country, err error)
	}{
		"success": {
			arrange: func() {
				arm.On("GetCountries", mock.Anything).Return([]*models.Country{country}, nil).Once()
			},
			assert: func(t *testing.T, actualCountries []*models.Country, err error) {
				require.NoError(t, err)
				require.NotNil(t, actualCountries)
				require.Len(t, actualCountries, 1)
				require.Equal(t, country, actualCountries[0])
			},
		},
		"failed": {
			arrange: func() {
				arm.On("GetCountries", mock.Anything).Return(([]*models.Country)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualCountries []*models.Country, err error) {
				require.Error(t, err)
				require.Zero(t, actualCountries)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			countries, err := addressService.GetCountries(context.Background())

			v.assert(t, countries, err)
		})
	}
}

func TestGetCityByID(t *testing.T) {
	tableTest := map[string]struct {
		cityID  uint
		arrange func()
		assert  func(t *testing.T, actualCity *models.City, err error)
	}{
		"success": {
			cityID: 1,
			arrange: func() {
				arm.On("GetCityByID", mock.Anything, uint(1)).Return(city, nil).Once()
			},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				require.NoError(t, err)
				require.NotZero(t, actualCity)
				require.Equal(t, city, actualCity)
			},
		},
		"failed": {
			cityID: 1,
			arrange: func() {
				arm.On("GetCityByID", mock.Anything, uint(1)).
					Return((*models.City)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				require.Error(t, err)
				require.Zero(t, actualCity)
			},
		},
		"empty city id": {
			arrange: func() {},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				var appErr *utilities.AppError
				require.Error(t, err)
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, utilities.ValidationError, appErr.Code)
				require.Equal(t, "city ID can't be empty", appErr.Message)
				require.Zero(t, actualCity)
				arm.AssertNotCalled(t, "GetCityByID")
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			actualCity, err := addressService.GetCityByID(context.Background(), v.cityID)

			v.assert(t, actualCity, err)
		})
	}
}

func TestCreateUserAddress(t *testing.T) {
	tableTest := map[string]struct{
		userID uint
		arrange func()
		assert func(t *testing.T, actual *models.Address, err error)
	}{
		"success": {
			userID: 1,
			arrange: func() {
				arm.On("CreateUserAddress", mock.Anything, &addressPayload, uint(1)).Return(address, nil).Once()
			},
			assert: func(t *testing.T, actual *models.Address, err error) {
				require.NoError(t, err)
				require.Equal(t, address, actual)
			},
		},
		"failed": {
			userID: 1,
			arrange: func() {
				arm.On("CreateUserAddress", mock.Anything, &addressPayload, uint(1)).
				Return((*models.Address)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actual *models.Address, err error) {
				require.Error(t, err)
				require.Zero(t, actual)
			},
		},
		"validation error": {
			arrange: func() {},
			assert: func(t *testing.T, actual *models.Address, err error) {
				var appErr *utilities.AppError
				require.Error(t, err)
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, utilities.ValidationError, appErr.Code)
				require.Zero(t, actual)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			actual, err := addressService.CreateUserAddress(context.Background(), &addressPayload, v.userID)

			v.assert(t, actual, err)
		})
	}
}