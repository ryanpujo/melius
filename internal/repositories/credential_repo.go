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
	FindByUsername(ctx context.Context, username string) (*models.User, error)
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

func (cr *CredentialRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT u.first_name, u.last_name, c.email, c.username, c.password
		FROM users u
		JOIN credentials c ON c.username = u.username
		WHERE u.username = $1
	`

	var user models.User

	err := cr.dB.QueryRowContext(ctx, query, username).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Credential.Email,
		&user.Credential.Username,
		&user.Credential.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a descriptive error if no user is found
			return nil, fmt.Errorf("user with username '%s' not found: %w", username, err)
		}
		// Return other database-related errors
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	return &user, nil
}
