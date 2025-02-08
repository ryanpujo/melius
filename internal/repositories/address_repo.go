package repositories

import (
	"context"
	"database/sql"
	"sort"
	"strings"

	"github.com/ryanpujo/melius/internal/models"
)

// AddressRepo defines the methods for managing address-related entities.
type AddressRepo interface {
	SaveCountry(ctx context.Context, country *models.CountryPayload, tx *sql.Tx) (*models.Country, error)
	SaveState(ctx context.Context, state *models.StatePayload, countryID uint, tx *sql.Tx) (*models.State, error)
	SaveCity(ctx context.Context, city *models.CityPayload, stateID uint, tx *sql.Tx) (*models.City, error)

	CreateUserAddress(ctx context.Context, address *models.AddressPayload, userID uint) (*models.Address, error)
	SaveAddress(ctx context.Context, address *models.AddressPayload, tx *sql.Tx) (*models.Address, error)

	GetCountries(ctx context.Context) ([]*models.Country, error)
	GetCityByID(ctx context.Context, cityID uint) (*models.City, error)
}

// addressRepo is the implementation of the AddressRepo interface.
type addressRepo struct {
	db *sql.DB
}

// NewAddressRepo creates a new instance of addressRepo.
func NewAddressRepo(db *sql.DB) AddressRepo {
	return &addressRepo{
		db: db,
	}
}

// saveEntity is a helper function to insert a record into the database and return its generated ID.
// Parameters:
//   - ctx: The context for managing request-scoped values.
//   - query: The SQL query string for inserting the record.
//   - args: The arguments for the SQL query.
//   - tx: An optional *sql.Tx transaction.
func (ar *addressRepo) saveEntity(ctx context.Context, query string, tx *sql.Tx, args ...interface{}) *sql.Row {
	var row *sql.Row

	if tx == nil {
		row = ar.db.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}

	return row
}

// SaveCountry inserts a new country into the "countries" table and returns its generated ID.
func (ar *addressRepo) SaveCountry(ctx context.Context, country *models.CountryPayload, tx *sql.Tx) (*models.Country, error) {
	query := `
		INSERT INTO countries (name) VALUES ($1) RETURNING id, name
	`

	var createdCountry models.Country
	row := ar.saveEntity(ctx, query, tx, country.Name)
	if err := row.Scan(&createdCountry.ID, &createdCountry.Name); err != nil {
		return nil, err
	}
	return &createdCountry, nil
}

func (ar *addressRepo) GetCountries(ctx context.Context) ([]*models.Country, error) {
	query := `
		SELECT
			c.id, c.name,
			s.id, s.name,
			ci.id, ci.name
		FROM countries c 
		JOIN states s ON c.id = s.country_id
		JOIN cities ci ON s.id = ci.state_id
		ORDER BY c.name, s.name, ci.name
	`

	rows, err := ar.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stateMap := make(map[uint]*models.State)
	countryMap := make(map[uint]*models.Country)
	countries := make([]*models.Country, 0, 30)

	for rows.Next() {
		var countryID uint
		var countryName string
		var stateID uint
		var stateName string
		var cityID uint
		var cityName string

		err := rows.Scan(&countryID, &countryName, &stateID, &stateName, &cityID, &cityName)
		if err != nil {
			return nil, err
		}

		country, exists := countryMap[countryID]
		if !exists {
			country = &models.Country{
				ID:     countryID,
				Name:   countryName,
				States: make([]*models.State, 0, 50),
			}
			countryMap[countryID] = country
		}

		state, exists := stateMap[stateID]
		if !exists {
			// If not, create a new state and add it to the map and country's slice
			state = &models.State{
				ID:     stateID,
				Name:   stateName,
				Cities: make([]*models.City, 0, 100),
			}
			stateMap[stateID] = state
			country.States = append(country.States, state)
		}

		// Add city to the corresponding state
		city := &models.City{
			ID:   cityID,
			Name: cityName,
		}
		state.Cities = append(state.Cities, city)
	}

	// Handle any error during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, v := range countryMap {
		countries = append(countries, v)
	}
	sort.Slice(countries, func(i, j int) bool {
		return strings.ToLower(countries[i].Name) < strings.ToLower(countries[j].Name)
	})

	return countries, nil
}

// SaveState inserts a new state into the "states" table and returns its generated ID.
func (ar *addressRepo) SaveState(ctx context.Context, state *models.StatePayload, countryID uint, tx *sql.Tx) (*models.State, error) {
	query := `
		INSERT INTO states (name, country_id) VALUES ($1, $2) RETURNING id, name
	`

	var createdState models.State
	row := ar.saveEntity(ctx, query, tx, state.Name, countryID)
	if err := row.Scan(&createdState.ID, &createdState.Name); err != nil {
		return nil, err
	}
	return &createdState, nil
}

// SaveCity inserts a new city into the "cities" table and returns its generated ID.
func (ar *addressRepo) SaveCity(ctx context.Context, city *models.CityPayload, stateID uint, tx *sql.Tx) (*models.City, error) {
	query := `
		INSERT INTO cities (name, state_id) VALUES ($1, $2) RETURNING id, name
	`
	var createdCity models.City
	row := ar.saveEntity(ctx, query, tx, city.Name, stateID)
	if err := row.Scan(&createdCity.ID, &createdCity.Name); err != nil {
		return nil, err
	}
	return &createdCity, nil
}

func (ar *addressRepo) GetCityByID(ctx context.Context, cityID uint) (*models.City, error) {
	query := `
		SELECT c.id, c.name, 
		s.id, s.name, 
		co.id, co.name 
		FROM cities c
		JOIN states s ON s.id = c.state_id
		JOIN countries co ON co.id = s.country_id
		WHERE c.id = $1
	`
	city := models.City{
		State: &models.State{
			Country: &models.Country{},
		},
	}
	err := ar.db.QueryRowContext(ctx, query, cityID).Scan(
		&city.ID,
		&city.Name,
		&city.State.ID,
		&city.State.Name,
		&city.State.Country.ID,
		&city.State.Country.Name,
	)
	if err != nil {
		return nil, err
	}

	return &city, nil
}

// SaveAddress inserts a new address into the database and returns its ID.
//
// The function takes a context for request-scoping, an Address model,
// a cityID to associate with the address, and an optional SQL transaction.
//
// Returns the ID of the new address or an error if the operation fails.
func (ar *addressRepo) SaveAddress(ctx context.Context, address *models.AddressPayload, tx *sql.Tx) (*models.Address, error) {
	query := `
		INSERT INTO addresses (address_line, postal_code, is_main, city_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, address_line, postal_code, is_main, city_id
	`
	createdAddress := models.Address{
		City: &models.City{},
	}
	row := ar.saveEntity(ctx, query, tx,
		address.AddressLine,
		address.PostalCode,
		address.IsMain,
		address.CityID,
	)
	err := row.Scan(
		&createdAddress.ID,
		&createdAddress.AddressLine,
		&createdAddress.PostalCode,
		&createdAddress.IsMain,
		&createdAddress.City.ID,
	)
	if err != nil {
		return nil, err
	}

	city, err := ar.GetCityByID(ctx, address.CityID)
	if err != nil {
		return nil, err
	}
	createdAddress.City = city

	return &createdAddress, nil
}

func (ar *addressRepo) CreateUserAddress(ctx context.Context, address *models.AddressPayload, userID uint) (*models.Address, error) {
	tx, err := ar.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO user_address (address_id, user_id) VALUES ($1, $2)
	`

	createdAddress, err := ar.SaveAddress(ctx, address, tx)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx, query, createdAddress.ID, userID)
	if err := row.Err(); err != nil {
		return nil, err
	}

	return createdAddress, tx.Commit()
}