package services

import (
	"context"
	"fmt"

	"github.com/ryanpujo/melius/internal/jwttoken"
	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type CredentialInterface interface {
	Write(ctx context.Context, payload models.UserPayload) (uint, error)
	FindByUsername(ctx context.Context, username string) (*models.Credential, error)
	Login(ctx context.Context, payload *models.LoginPayload) (string, error)
}

type CredentialService struct {
	credRepo repositories.CredentialInterface
}

func NewCredService(credRepo repositories.CredentialInterface) *CredentialService {
	return &CredentialService{
		credRepo: credRepo,
	}
}

var HashPassword = func(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (cs *CredentialService) Write(ctx context.Context, payload models.UserPayload) (uint, error) {
	password, err := HashPassword(payload.CredentialPayload.Password)
	if err != nil {
		return 0, err
	}
	payload.CredentialPayload.Password = string(password)
	return cs.credRepo.Write(ctx, payload)
}

func (cs *CredentialService) FindByUsername(ctx context.Context, username string) (*models.Credential, error) {

	return cs.credRepo.FindByUsername(ctx, username)
}

var CompareHashAndPassword = func(hash string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func (cs *CredentialService) Login(ctx context.Context, payload *models.LoginPayload) (string, error) {
	credFound, err := cs.FindByUsername(ctx, payload.Username)
	if err != nil {
		return "", fmt.Errorf("credential not found: %w", err)
	}

	err = CompareHashAndPassword(credFound.Password, payload.Password)
	if err != nil {
		return "", fmt.Errorf("credential not found: %w", err)
	}

	return jwttoken.GenerateJWT(credFound.Username)
}
