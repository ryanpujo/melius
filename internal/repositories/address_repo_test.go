package repositories_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/repositories"
	"github.com/stretchr/testify/require"
)

var (
	countryPayload = &models.CountryPayload{
		Name: "Indonesia",
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
		Name: "Jakarta Timur",
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
	query := regexp.QuoteMeta(fmt.Sprintf(repositories.PreparedInsertQuery, "countries", "name", "name"))
	tableTest := map[string]struct {
		tx      func() *sql.Tx
		arrange func()
		assert  func(t *testing.T, actualCountry *models.Country, err error)
	}{
		"success with no tx": {
			tx: func() *sql.Tx {
				return nil
			},
			arrange: func() {
				row := mock.NewRows([]string{"id", "name"}).AddRow(country.ID, country.Name)
				mock.ExpectQuery(query).WithArgs(countryPayload.Name).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualCountry *models.Country, err error) {
				require.NoError(t, err)
				require.Equal(t, country, actualCountry)
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
				row := mock.NewRows([]string{"id", "name"}).AddRow(country.ID, country.Name)
				mock.ExpectQuery(query).WithArgs(countryPayload.Name).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualCountry *models.Country, err error) {
				require.NoError(t, err)
				require.Equal(t, country, actualCountry)
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
				mock.ExpectQuery(query).WithArgs(countryPayload.Name).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualCountry *models.Country, err error) {
				require.Error(t, err)
				require.Zero(t, actualCountry)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressRepo.SaveCountry(context.Background(), countryPayload, v.tx())

			v.assert(t, id, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestSaveState(t *testing.T) {
	query := regexp.QuoteMeta(fmt.Sprintf(repositories.PreparedInsertQuery, "states", "name, country_id", "name, country_id"))
	tableTest := map[string]struct {
		tx      func() *sql.Tx
		arrange func()
		assert  func(t *testing.T, actualState *models.State, err error)
	}{
		"success with no tx": {
			tx: func() *sql.Tx {
				return nil
			},
			arrange: func() {
				row := mock.NewRows([]string{"id", "name"}).AddRow(state.ID, state.Name)
				mock.ExpectQuery(query).WithArgs(statePayload.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualState *models.State, err error) {
				require.NoError(t, err)
				require.Equal(t, state, actualState)
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
				row := mock.NewRows([]string{"id", "name"}).AddRow(state.ID, state.Name)
				mock.ExpectQuery(query).WithArgs(statePayload.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualState *models.State, err error) {
				require.NoError(t, err)
				require.Equal(t, actualState, actualState)
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
				mock.ExpectQuery(query).WithArgs(statePayload.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualState *models.State, err error) {
				require.Error(t, err)
				require.Zero(t, actualState)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressRepo.SaveState(context.Background(), &statePayload, 2, v.tx())

			v.assert(t, id, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestSaveCity(t *testing.T) {
	query := regexp.QuoteMeta(fmt.Sprintf(repositories.PreparedInsertQuery, "cities", "name, state_id", "name, state_id"))
	tableTest := map[string]struct {
		tx      func() *sql.Tx
		arrange func()
		assert  func(t *testing.T, actualCity *models.City, err error)
	}{
		"success with no tx": {
			tx: func() *sql.Tx {
				return nil
			},
			arrange: func() {
				row := mock.NewRows([]string{"id", "name"}).AddRow(city.ID, city.Name)
				mock.ExpectQuery(query).WithArgs(cityPayload.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				require.NoError(t, err)
				require.Equal(t, city, actualCity)
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
				row := mock.NewRows([]string{"id", "name"}).AddRow(city.ID, city.Name)
				mock.ExpectQuery(query).WithArgs(cityPayload.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				require.NoError(t, err)
				require.Equal(t, city, actualCity)
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
				mock.ExpectQuery(query).WithArgs(cityPayload.Name, 2).
					WillReturnRows(row)
			},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				require.Error(t, err)
				require.Zero(t, actualCity)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressRepo.SaveCity(context.Background(), &cityPayload, 2, v.tx())

			v.assert(t, id, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func getCityByIDQuery(cityID int) *sqlmock.ExpectedQuery {
	row := mock.NewRows([]string{"id", "name", "id", "name", "id", "name"}).
		AddRow(
			city.ID, city.Name,
			state.ID, state.Name,
			country.ID, country.Name,
		)
	query := `
		SELECT c.id, c.name, 
		s.id, s.name, 
		co.id, co.name 
		FROM cities c
		JOIN states s ON s.id = c.state_id
		JOIN countries co ON co.id = s.country_id
		WHERE c.id = $1
	`
	query = regexp.QuoteMeta(query)
	return mock.ExpectQuery(query).WithArgs(cityID).WillReturnRows(row)
}

func TestGetCityByID(t *testing.T) {
	state = &models.State{
		ID:      1,
		Name:    "Jakarta",
		Country: country,
	}
	city := &models.City{
		ID:    1,
		Name:  "Jakarta Timur",
		State: state,
	}
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualCity *models.City, err error)
	}{
		"success": {
			arrange: func() {
				getCityByIDQuery(1)
			},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				require.NoError(t, err)
				require.Equal(t, city, actualCity)
			},
		},
		"scan error": {
			arrange: func() {
				getCityByIDQuery(1).WillReturnError(errors.New("failed"))
			},
			assert: func(t *testing.T, actualCity *models.City, err error) {
				require.Error(t, err)
				require.Zero(t, actualCity)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			actual, err := addressRepo.GetCityByID(context.Background(), 1)

			v.assert(t, actual, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func saveAddressMockExpectation(addrTest *models.Address) *sqlmock.ExpectedQuery {
	query := `
		WITH created_address AS (
			INSERT INTO addresses (address_line, postal_code, is_main, city_id) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id, address_line, postal_code, is_main, city_id
		)
		SELECT 
			ca.id, ca.address_line, ca.postal_code, ca.is_main,
			c.id, c.name,
			s.id, s.name,
			co.id, co.name
		FROM created_address ca
		JOIN cities c ON c.id = ca.city_id
		JOIN states s ON s.id = c.state_id
		JOIN countries co ON co.id = s.country_id
	`
	query = regexp.QuoteMeta(query)
	columns := []string{
		"id", "address_line", "postal_code", "is_main",
		"city_id", "city_name",
		"state_id", "state_name",
		"country_id", "country_name",
	}

	// Create a row with the expected data.
	rows := sqlmock.NewRows(columns).
		AddRow(addrTest.ID, addrTest.AddressLine, addrTest.PostalCode, addrTest.IsMain,
			addrTest.City.ID, addrTest.City.Name,
			addrTest.City.State.ID, addrTest.City.State.Name,
			addrTest.City.State.Country.ID, addrTest.City.State.Country.Name)

	return mock.ExpectQuery(query).WithArgs(
		addressPayload.AddressLine,
		addressPayload.PostalCode,
		addressPayload.IsMain,
		address.City.ID,
	).
		WillReturnRows(rows)
}

func TestSaveAddress(t *testing.T) {
	state = &models.State{
		ID:      1,
		Name:    "Jakarta",
		Country: country,
	}
	city := &models.City{
		ID:    1,
		Name:  "Jakarta Timur",
		State: state,
	}
	address = &models.Address{
		ID:          1,
		AddressLine: "jl. mayjen sutoyo kel. cawang kec kramat jati rt.007/011",
		PostalCode:  "12630",
		IsMain:      true,
		City:        city,
	}
	tableTest := map[string]struct {
		tx      func() *sql.Tx
		arrange func()
		assert  func(t *testing.T, actualAddress *models.Address, err error)
	}{
		"success with no tx": {
			tx: func() *sql.Tx {
				return nil
			},
			arrange: func() {
				saveAddressMockExpectation(address)
			},
			assert: func(t *testing.T, actualAddress *models.Address, err error) {
				require.NoError(t, err)
				require.Equal(t, address, actualAddress)
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
				saveAddressMockExpectation(address)
			},
			assert: func(t *testing.T, actualAddress *models.Address, err error) {
				require.NoError(t, err)
				require.Equal(t, address, actualAddress)
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
				saveAddressMockExpectation(address).WillReturnError(errors.New("failed"))
			},
			assert: func(t *testing.T, actualAddress *models.Address, err error) {
				require.Error(t, err)
				require.Zero(t, actualAddress)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			id, err := addressRepo.SaveAddress(context.Background(), &addressPayload, v.tx())

			v.assert(t, id, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetCountries(t *testing.T) {
	expectedQuery := regexp.QuoteMeta(`
		SELECT
			c.id, c.name,
			s.id, s.name,
			ci.id, ci.name
		FROM countries c 
		JOIN states s ON c.id = s.country_id
		JOIN cities ci ON s.id = ci.state_id
		ORDER BY c.name, s.name, ci.name
	`)

	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualCountries []*models.Country, err error)
	}{
		"success": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "id", "name", "id", "name"}).
					AddRow(1, "USA", 10, "California", 100, "Los Angeles").
					AddRow(1, "USA", 10, "California", 101, "San Francisco").
					AddRow(2, "Canada", 20, "Ontario", 200, "Toronto")

				// Expect the query to be executed and return our rows.
				mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualCountries []*models.Country, err error) {
				require.NoError(t, err)
				require.NotNil(t, actualCountries)
				require.Len(t, actualCountries, 2)
				require.Len(t, actualCountries[1].States[0].Cities, 2)
				require.Equal(t, "Canada", actualCountries[0].Name)
			},
		},
		"scan error": {
			arrange: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "id", "name", "id", "name"}).
					AddRow(1, "USA", 10, "California", "sffe", "Los Angeles").
					AddRow(1, "USA", 10, "California", 101, "San Francisco").
					AddRow(2, "Canada", 20, "Ontario", 200, "Toronto")

				// Expect the query to be executed and return our rows.
				mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualCountries []*models.Country, err error) {
				require.Error(t, err)
				require.Nil(t, actualCountries)
			},
		},
		"query error": {
			arrange: func() {
				mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("failed"))
			},
			assert: func(t *testing.T, actualCountries []*models.Country, err error) {
				require.Error(t, err)
				require.Nil(t, actualCountries)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			countries, err := addressRepo.GetCountries(context.Background())

			v.assert(t, countries, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCreateUserAddress(t *testing.T) {
	state = &models.State{
		ID:      1,
		Name:    "Jakarta",
		Country: country,
	}
	city := &models.City{
		ID:    1,
		Name:  "Jakarta Timur",
		State: state,
	}
	address = &models.Address{
		ID:          1,
		AddressLine: "jl. mayjen sutoyo kel. cawang kec kramat jati rt.007/011",
		PostalCode:  "12630",
		IsMain:      true,
		City:        city,
	}
	expectedQuery := regexp.QuoteMeta(`
		INSERT INTO user_address (address_id, user_id) VALUES ($1, $2)
	`)
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actual *models.Address, err error)
	}{
		"success": {
			arrange: func() {
				mock.ExpectBegin()

				saveAddressMockExpectation(address)
				mock.ExpectExec(expectedQuery).WithArgs(1, 2).WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()
			},
			assert: func(t *testing.T, actual *models.Address, err error) {
				require.NoError(t, err)
				require.NotZero(t, actual)
				require.Equal(t, address, actual)
			},
		},
		"failed to start tx": {
			arrange: func() {
				mock.ExpectBegin().WillReturnError(errors.New("failed"))
			},
			assert: func(t *testing.T, actual *models.Address, err error) {
				require.Error(t, err)
				require.Zero(t, actual)
			},
		},
		"failed to save address": {
			arrange: func() {
				mock.ExpectBegin()
				saveAddressMockExpectation(address).WillReturnError(errors.New("failed"))
			},
			assert: func(t *testing.T, actual *models.Address, err error) {
				require.Error(t, err)
				require.Zero(t, actual)
			},
		},
		"row error": {
			arrange: func() {
				mock.ExpectBegin()

				saveAddressMockExpectation(address)
				mock.ExpectExec(expectedQuery).WithArgs(1, 2).WillReturnError(errors.New("failed"))
			},
			assert: func(t *testing.T, actual *models.Address, err error) {
				require.Error(t, err)
				require.Zero(t, actual)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			actual, err := addressRepo.CreateUserAddress(context.Background(), &addressPayload, 2)

			v.assert(t, actual, err)
		})
	}
	err := mock.ExpectationsWereMet()
	require.NoError(t, err)
}
