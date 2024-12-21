package services

import (
	"context"

	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type CredentialInterface interface {
	Write(ctx context.Context, payload models.UserPayload) (uint, error)
}

type CredentialService struct {
	credRepo repositories.CredentialInterface
}

func NewCredService(credRepo repositories.CredentialInterface) *CredentialService {
	return &CredentialService{
		credRepo: credRepo,
	}
}

func (cs *CredentialService) Write(ctx context.Context, payload models.UserPayload) (uint, error) {
	password, err :=
		bcrypt.GenerateFromPassword([]byte(payload.CredentialPayload.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	payload.CredentialPayload.Password = string(password)
	return cs.credRepo.Write(ctx, payload)
}
