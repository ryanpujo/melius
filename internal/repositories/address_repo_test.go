package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ryanpujo/melius/internal/models"
	"github.com/stretchr/testify/require"
)

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
		tx      func() *sql.Tx
		arrange func()
		assert  func(t *testing.T, actualID uint, err error)
	}{
		"success with no tx": {
			tx: func() *sql.Tx {
				return nil
			},
			arrange: func() {
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO countries").WithArgs(country.Name).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
			},
		},
		"succes with tx": {
			tx: func() *sql.Tx {
				tx, err := db.Begin()
				require.NoError(t, err)
				return tx
			},
			arrange: func() {
				mock.ExpectBegin()
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO countries").WithArgs(country.Name).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
			},
		},
		"row error": {
			tx: func() *sql.Tx {
				tx, err := db.Begin()
				require.NoError(t, err)
				return tx
			},
			arrange: func() {
				mock.ExpectBegin()
				row := mock.NewRows([]string{"id"}).AddRow("string")
				mock.ExpectQuery("INSERT INTO countries").WithArgs(country.Name).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressRepo.SaveCountry(context.Background(), country, v.tx())

			v.assert(t, id, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestSaveState(t *testing.T) {
	tableTest := map[string]struct {
		tx      func() *sql.Tx
		arrange func()
		assert  func(t *testing.T, actualID uint, err error)
	}{
		"success with no tx": {
			tx: func() *sql.Tx {
				return nil
			},
			arrange: func() {
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO states").WithArgs(state.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
			},
		},
		"succes with tx": {
			tx: func() *sql.Tx {
				tx, err := db.Begin()
				require.NoError(t, err)
				return tx
			},
			arrange: func() {
				mock.ExpectBegin()
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO states").WithArgs(state.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
			},
		},
		"row error": {
			tx: func() *sql.Tx {
				tx, err := db.Begin()
				require.NoError(t, err)
				return tx
			},
			arrange: func() {
				mock.ExpectBegin()
				row := mock.NewRows([]string{"id"}).AddRow("string")
				mock.ExpectQuery("INSERT INTO states").WithArgs(state.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressRepo.SaveState(context.Background(), state, 2, v.tx())

			v.assert(t, id, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestSaveCity(t *testing.T) {
	tableTest := map[string]struct {
		tx      func() *sql.Tx
		arrange func()
		assert  func(t *testing.T, actualID uint, err error)
	}{
		"success with no tx": {
			tx: func() *sql.Tx {
				return nil
			},
			arrange: func() {
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO cities").WithArgs(city.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
			},
		},
		"succes with tx": {
			tx: func() *sql.Tx {
				tx, err := db.Begin()
				require.NoError(t, err)
				return tx
			},
			arrange: func() {
				mock.ExpectBegin()
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO cities").WithArgs(city.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
			},
		},
		"row error": {
			tx: func() *sql.Tx {
				tx, err := db.Begin()
				require.NoError(t, err)
				return tx
			},
			arrange: func() {
				mock.ExpectBegin()
				row := mock.NewRows([]string{"id"}).AddRow("string")
				mock.ExpectQuery("INSERT INTO cities").WithArgs(city.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressRepo.SaveCity(context.Background(), city, 2, v.tx())

			v.assert(t, id, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestSaveAddress(t *testing.T) {
	tableTest := map[string]struct {
		tx      func() *sql.Tx
		arrange func()
		assert  func(t *testing.T, actualID uint, err error)
	}{
		"success with no tx": {
			tx: func() *sql.Tx {
				return nil
			},
			arrange: func() {
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO addresses").WithArgs(
					address.AddressLine,
					address.PostalCode,
					address.IsMain,
					2,
				).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
			},
		},
		"succes with tx": {
			tx: func() *sql.Tx {
				tx, err := db.Begin()
				require.NoError(t, err)
				return tx
			},
			arrange: func() {
				mock.ExpectBegin()
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO addresses").WithArgs(
					address.AddressLine,
					address.PostalCode,
					address.IsMain,
					2,
				).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), actualID)
			},
		},
		"row error": {
			tx: func() *sql.Tx {
				tx, err := db.Begin()
				require.NoError(t, err)
				return tx
			},
			arrange: func() {
				mock.ExpectBegin()
				row := mock.NewRows([]string{"id"}).AddRow("string")
				mock.ExpectQuery("INSERT INTO addresses").WithArgs(
					address.AddressLine,
					address.PostalCode,
					address.IsMain,
					2,
				).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualID uint, err error) {
				require.Error(t, err)
				require.Zero(t, actualID)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressRepo.SaveAddress(context.Background(), address, 2, v.tx())

			v.assert(t, id, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}
