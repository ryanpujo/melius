package services

import (
	"context"

	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/repositories"
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
	return cs.credRepo.Write(ctx, payload)
}
