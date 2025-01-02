package repositories

import (
	"context"
	"database/sql"

	"github.com/ryanpujo/melius/internal/models"
)

// AddressRepo defines the methods for managing address-related entities.
type AddressRepo interface {
	SaveCountry(ctx context.Context, country models.Country, tx *sql.Tx) (uint, error)
	SaveState(ctx context.Context, state models.State, countryID uint, tx *sql.Tx) (uint, error)
	SaveCity(ctx context.Context, city models.City, stateID uint, tx *sql.Tx) (uint, error)
	SaveAddress(ctx context.Context, address models.Address, cityID uint, tx *sql.Tx) (uint, error)
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
func (ar *addressRepo) saveEntity(ctx context.Context, query string, tx *sql.Tx, args ...interface{}) (uint, error) {
	var id uint
	var row *sql.Row

	if tx == nil {
		row = ar.db.QueryRowContext(ctx, query, args...)
	} else {
		row = tx.QueryRowContext(ctx, query, args...)
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// SaveCountry inserts a new country into the "countries" table and returns its generated ID.
func (ar *addressRepo) SaveCountry(ctx context.Context, country models.Country, tx *sql.Tx) (uint, error) {
	query := `
		INSERT INTO countries (name) VALUES ($1) RETURNING id
	`
	return ar.saveEntity(ctx, query, tx, country.Name)
}

// SaveState inserts a new state into the "states" table and returns its generated ID.
func (ar *addressRepo) SaveState(ctx context.Context, state models.State, countryID uint, tx *sql.Tx) (uint, error) {
	query := `
		INSERT INTO states (name, country_id) VALUES ($1, $2) RETURNING id
	`
	return ar.saveEntity(ctx, query, tx, state.Name, countryID)
}

// SaveCity inserts a new city into the "cities" table and returns its generated ID.
func (ar *addressRepo) SaveCity(ctx context.Context, city models.City, stateID uint, tx *sql.Tx) (uint, error) {
	query := `
		INSERT INTO cities (name, state_id) VALUES ($1, $2) RETURNING id
	`
	return ar.saveEntity(ctx, query, tx, city.Name, stateID)
}

// SaveAddress inserts a new address into the database and returns its ID.
//
// The function takes a context for request-scoping, an Address model,
// a cityID to associate with the address, and an optional SQL transaction.
//
// Returns the ID of the new address or an error if the operation fails.
func (ar *addressRepo) SaveAddress(ctx context.Context, address models.Address, cityID uint, tx *sql.Tx) (uint, error) {
	query := `
		INSERT INTO addresses (address_line, postal_code, is_main, city_id) 
		VALUES ($1, $2, $3, $4) RETURNING id
	`
	return ar.saveEntity(ctx, query, tx,
		address.AddressLine,
		address.PostalCode,
		address.IsMain,
		cityID,
	)
}
