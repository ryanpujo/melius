package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ryanpujo/melius/internal/models"
)

type CredentialInterface interface {
	Write(ctx context.Context, payload models.UserPayload) (uint, error)
	FindByUsername(ctx context.Context, username string) (*models.Credential, error)
}

type CredentialRepo struct {
	dB *sql.DB
}

func NewCredentialRepo(db *sql.DB) *CredentialRepo {
	return &CredentialRepo{
		dB: db,
	}
}

// Write inserts a user's credentials and basic information into the database.
// It performs the operation in a transaction to ensure atomicity.
// Parameters:
//   - ctx: The context for the operation, allowing for cancellation and timeout.
//   - payload: The user's data to insert.
//
// Returns:
//   - The generated user ID on success.
//   - An error if the operation fails.
func (cr *CredentialRepo) Write(ctx context.Context, payload models.UserPayload) (uint, error) {
	userQuery := `
		INSERT INTO users (first_name, last_name, username, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	credentialQuery := `
		INSERT INTO credentials (email, username, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING username
	`

	tx, err := cr.dB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var username string

	err = tx.QueryRowContext(ctx, credentialQuery,
		payload.CredentialPayload.Email,
		payload.CredentialPayload.Username,
		payload.CredentialPayload.Password,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
	).Scan(&username)
	if err != nil {
		return 0, err
	}

	var id uint

	err = tx.QueryRowContext(ctx, userQuery,
		payload.FirstName,
		payload.LastName,
		username,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

// FindByUsername retrieves a user's credentials from the database by username.
// It returns the corresponding credential or an error if the user is not found.
func (cr *CredentialRepo) FindByUsername(ctx context.Context, username string) (*models.Credential, error) {
	query := `
		SELECT email, username, password 
		FROM credentials
		WHERE username = $1
	`

	// Execute the query and scan the result into the credential struct
	row := cr.dB.QueryRowContext(ctx, query, username)
	var credential models.Credential
	if err := row.Scan(
		&credential.Email,
		&credential.Username,
		&credential.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			// Return a descriptive error if no user is found
			return nil, fmt.Errorf("user with username '%s' not found: %w", username, err)
		}
		// Return other database-related errors
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Return the populated credential struct
	return &credential, nil
}
