package services_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type CredRepoMock struct {
	mock.Mock
}

func (crm *CredRepoMock) Write(ctx context.Context, payload models.UserPayload) (uint, error) {
	args := crm.Called(ctx, payload)
	return uint(args.Int(0)), args.Error(1)
}

var (
	credService services.CredentialService
	crm         *CredRepoMock
)

func TestMain(m *testing.M) {
	crm = new(CredRepoMock)
	credService = *services.NewCredService(crm)
	os.Exit(m.Run())
}

func TestWriteUser(t *testing.T) {
	credential := models.CredentialPayload{
		Email:    "ryanpujo@gmail.com",
		Username: "ryanpujo",
		Password: "okeoke",
	}
	user := models.UserPayload{
		FirstName:         "Ryan",
		LastName:          "Pujo",
		CredentialPayload: credential,
	}
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, id uint, err error)
	}{
		"success": {
			arrange: func() {
				crm.On("Write", mock.Anything, user).Return(1, nil).Once()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), id)
			},
		},
		"failed": {
			arrange: func() {
				crm.On("Write", mock.Anything, user).Return(0, errors.New("failed")).Once()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Error(t, err)
				require.Zero(t, id)
			},
		},
	}

	for key, v := range tableTest {
		t.Run(key, func(t *testing.T) {
			v.arrange()

			id, err := credService.Write(context.Background(), user)

			v.assert(t, id, err)
		})
	}
}
