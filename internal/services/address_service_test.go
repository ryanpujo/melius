package services_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/ryanpujo/melius/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type addressRepoMock struct {
	mock.Mock
}

func (ar *addressRepoMock) SaveCountry(ctx context.Context, country models.Country, tx *sql.Tx) (uint, error) {
	args := ar.Called(ctx, country, tx)
	return uint(args.Int(0)), args.Error(1)
}

func (ar *addressRepoMock) SaveState(ctx context.Context, state models.State, countryID uint, tx *sql.Tx) (uint, error) {
	args := ar.Called(ctx, state, countryID, tx)
	return uint(args.Int(0)), args.Error(1)
}

func (ar *addressRepoMock) SaveCity(ctx context.Context, city models.City, stateID uint, tx *sql.Tx) (uint, error) {
	args := ar.Called(ctx, city, stateID, tx)
	return uint(args.Int(0)), args.Error(1)
}

func (ar *addressRepoMock) SaveAddress(ctx context.Context, address models.Address, cityID uint, tx *sql.Tx) (uint, error) {
	args := ar.Called(ctx, address, cityID, tx)
	return uint(args.Int(0)), args.Error(1)
}

var (
	country = models.Country{
		Name: "Indonesia",
	}
	state = models.State{
		Name:    "Jakarta",
		Country: country,
	}
	city = models.City{
		Name:  "Jakarta Timur",
		State: state,
	}
	address = models.Address{
		AddressLine: "jl. mayjen sutoyo kel. cawang kec kramat jati rt.007/011",
		PostalCode:  "12630",
		IsMain:      true,
		City:        city,
	}
)

func TestSaveCountry(t *testing.T) {
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualID uint, err error)
	}{
		"succes": {
			arrange: func() {
				arm.On("SaveCountry", mock.Anything, country, (*sql.Tx)(nil)).Return(1, nil).Once()
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
				arm.AssertCalled(t, "SaveCountry", context.Background(), country, (*sql.Tx)(nil))
			},
		},
		"failed": {
			arrange: func() {
				arm.On("SaveCountry", mock.Anything, country, (*sql.Tx)(nil)).Return(0, errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
				arm.AssertCalled(t, "SaveCountry", context.Background(), country, (*sql.Tx)(nil))
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressService.SaveCountry(context.Background(), country)

			v.assert(t, id, err)
		})
	}
}

func TestSaveState(t *testing.T) {
	tableTest := map[string]struct {
		countryID uint
		arrange   func()
		assert    func(t *testing.T, actualID uint, err error)
	}{
		"success": {
			countryID: 1,
			arrange: func() {
				arm.On("SaveState", mock.Anything, state, uint(1), (*sql.Tx)(nil)).Return(1, nil).Once()
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
				arm.AssertCalled(t, "SaveState", context.Background(), state, uint(1), (*sql.Tx)(nil))
			},
		},
		"failed": {
			countryID: 1,
			arrange: func() {
				arm.On("SaveState", mock.Anything, state, uint(1), (*sql.Tx)(nil)).Return(0, errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
				arm.AssertCalled(t, "SaveState", context.Background(), state, uint(1), (*sql.Tx)(nil))
			},
		},
		"empty country id": {
			countryID: 0,
			arrange:   func() {},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
				require.Equal(t, "country ID cannot be empty", err.Error())
				arm.AssertNotCalled(t, "SaveState")
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressService.SaveState(context.Background(), state, v.countryID)

			v.assert(t, id, err)
		})
	}
}

func TestSaveCity(t *testing.T) {
	tableTest := map[string]struct {
		stateID uint
		arrange func()
		assert  func(t *testing.T, actualID uint, err error)
	}{
		"success": {
			stateID: 1,
			arrange: func() {
				arm.On("SaveCity", mock.Anything, city, uint(1), (*sql.Tx)(nil)).Return(1, nil).Once()
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
				arm.AssertCalled(t, "SaveCity", context.Background(), city, uint(1), (*sql.Tx)(nil))
			},
		},
		"failed": {
			stateID: 1,
			arrange: func() {
				arm.On("SaveCity", mock.Anything, city, uint(1), (*sql.Tx)(nil)).Return(0, errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
				arm.AssertCalled(t, "SaveCity", context.Background(), city, uint(1), (*sql.Tx)(nil))
			},
		},
		"empty country id": {
			stateID: 0,
			arrange: func() {},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
				require.Equal(t, "state ID cannot be empty", err.Error())
				arm.AssertNotCalled(t, "SaveCity")
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressService.SaveCity(context.Background(), city, v.stateID)

			v.assert(t, id, err)
		})
	}
}

func TestSaveAddress(t *testing.T) {
	tableTest := map[string]struct {
		cityID  uint
		arrange func()
		assert  func(t *testing.T, actualID uint, err error)
	}{
		"success": {
			cityID: 1,
			arrange: func() {
				arm.On("SaveAddress", mock.Anything, address, uint(1), (*sql.Tx)(nil)).Return(1, nil).Once()
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
				arm.AssertCalled(t, "SaveAddress", context.Background(), address, uint(1), (*sql.Tx)(nil))
			},
		},
		"failed": {
			cityID: 1,
			arrange: func() {
				arm.On("SaveAddress", mock.Anything, address, uint(1), (*sql.Tx)(nil)).Return(0, errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
				arm.AssertCalled(t, "SaveAddress", context.Background(), address, uint(1), (*sql.Tx)(nil))
			},
		},
		"empty country id": {
			cityID:  0,
			arrange: func() {},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
				require.Equal(t, "city ID cannot be empty", err.Error())
				arm.AssertNotCalled(t, "SaveAddress")
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressService.SaveAddress(context.Background(), address, v.cityID)

			v.assert(t, id, err)
		})
	}
}
