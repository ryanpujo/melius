package services

import (
	"context"
	"fmt"

	"github.com/ryanpujo/melius/internal/jwttoken"
	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// CredentialInterface defines the contract for credential-related operations.
type CredentialInterface interface {
	Write(ctx context.Context, payload models.UserPayload) (uint, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Login(ctx context.Context, payload *models.LoginPayload) (string, error)
}

// CredentialService implements the CredentialInterface and provides business logic.
type CredentialService struct {
	credRepo repositories.CredentialInterface
}

// NewCredentialService creates a new instance of CredentialService.
func NewCredentialService(credRepo repositories.CredentialInterface) *CredentialService {
	return &CredentialService{
		credRepo: credRepo,
	}
}

// HashPassword generates a bcrypt hash of the given password.
var HashPassword = func(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CompareHashAndPassword verifies that a plain-text password matches a bcrypt hash.
var CompareHashAndPassword = func(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

// Write creates a new user credential and stores it in the repository.
// It hashes the password before saving.
func (cs *CredentialService) Write(ctx context.Context, payload models.UserPayload) (uint, error) {
	passwordHash, err := HashPassword(payload.CredentialPayload.Password)
	if err != nil {
		return 0, err
	}

	// Replace the plain-text password with the hashed password.
	payload.CredentialPayload.Password = passwordHash

	// Delegate the write operation to the repository.
	return cs.credRepo.Write(ctx, payload)
}

// FindByUsername retrieves a credential by username from the repository.
func (cs *CredentialService) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	credential, err := cs.credRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to find credential: %w", err)
	}
	return credential, nil
}

// Login authenticates a user by username and password, returning a JWT if successful.
func (cs *CredentialService) Login(ctx context.Context, payload *models.LoginPayload) (string, error) {
	// Retrieve the credential by username.
	user, err := cs.FindByUsername(ctx, payload.Username)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	// Verify the provided password matches the stored hash.
	if err := CompareHashAndPassword(user.Credential.Password, payload.Password); err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	// Generate a JWT token for the authenticated user.
	token, err := jwttoken.GenerateJWT(user.Credential.Username)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	return token, nil
}

// Documentation Summary:
// 1. CredentialInterface: Defines the contract for credential operations.
// 2. CredentialService: Implements the business logic for credential operations.
// 3. HashPassword: Hashes a plain-text password using bcrypt.
// 4. CompareHashAndPassword: Verifies a password against a bcrypt hash.
// 5. Write: Handles the creation of new credentials with password hashing.
// 6. FindByUsername: Retrieves credentials by username.
// 7. Login: Authenticates a user and generates a JWT on successful login.
