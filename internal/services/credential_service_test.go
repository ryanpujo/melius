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

func (crm *CredRepoMock) FindByUsername(ctx context.Context, username string) (*models.Credential, error) {
	args := crm.Called(ctx, username)
	return args.Get(0).(*models.Credential), args.Error(1)
}

var (
	credService       services.CredentialService
	crm               *CredRepoMock
	hashFunc          = services.HashPassword
	compareFunc       = services.CompareHashAndPassword
	credentialPayload = models.CredentialPayload{
		Email:    "ryanpujo@gmail.com",
		Username: "ryanpujo",
		Password: "okeoke",
	}

	credential = models.Credential{
		Email:    "ryanpujo@gmail.com",
		Username: "ryanpujo",
		Password: "okeoke",
	}
	userPayload = models.UserPayload{
		FirstName:         "Ryan",
		LastName:          "Pujo",
		CredentialPayload: credentialPayload,
	}
)

func TestMain(m *testing.M) {
	crm = new(CredRepoMock)
	credService = *services.NewCredService(crm)
	os.Exit(m.Run())
}

func TestWriteUser(t *testing.T) {
	tableTest := map[string]struct {
		arrange  func()
		assert   func(t *testing.T, id uint, err error)
		teardown func()
	}{
		"success": {
			arrange: func() {
				crm.On("Write", mock.Anything, mock.Anything).Return(1, nil).Once()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), id)
			},
			teardown: func() {},
		},
		"failed": {
			arrange: func() {
				crm.On("Write", mock.Anything, mock.Anything).Return(0, errors.New("failed")).Once()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Error(t, err)
				require.Zero(t, id)
			},
			teardown: func() {},
		},
		"bcrypt failed": {
			arrange: func() {
				services.HashPassword = func(password string) (string, error) {
					return "", errors.New("failed to hash")
				}
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Error(t, err)
				require.Zero(t, id)
			},
			teardown: func() {
				services.HashPassword = hashFunc
			},
		},
	}

	for key, v := range tableTest {
		t.Run(key, func(t *testing.T) {
			v.arrange()

			id, err := credService.Write(context.Background(), userPayload)

			v.assert(t, id, err)

			v.teardown()
		})
	}
}

func TestFindByUsername(t *testing.T) {
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, cred *models.Credential, err error)
	}{
		"success": {
			arrange: func() {
				crm.On("FindByUsername", mock.Anything, mock.Anything).Return(&credential, nil).Once()
			},
			assert: func(t *testing.T, cred *models.Credential, err error) {
				require.NoError(t, err)
				require.Equal(t, &credential, cred)
				require.True(t, crm.AssertCalled(t, "FindByUsername", context.Background(), "ryanpujo"))
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			cred, err := credService.FindByUsername(context.Background(), "ryanpujo")

			v.assert(t, cred, err)
		})
	}
}

func TestLogin(t *testing.T) {
	tableTest := map[string]struct {
		arrange  func()
		assert   func(t *testing.T, jwt string, err error)
		teardown func()
	}{
		"success": {
			arrange: func() {
				crm.On("FindByUsername", mock.Anything, mock.Anything).Return(&credential, nil).Once()
				services.CompareHashAndPassword = func(hash, plain string) error {
					return nil
				}
			},
			assert: func(t *testing.T, jwt string, err error) {
				require.NoError(t, err)
				require.NotZero(t, jwt)
			},
			teardown: func() {
				services.CompareHashAndPassword = compareFunc
			},
		},
		"credential not found": {
			arrange: func() {
				crm.On("FindByUsername", mock.Anything, mock.Anything).
					Return((*models.Credential)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, jwt string, err error) {
				require.Error(t, err)
				require.Zero(t, jwt)
			},
			teardown: func() {},
		},
		"wrong password": {
			arrange: func() {
				crm.On("FindByUsername", mock.Anything, mock.Anything).Return(&credential, nil).Once()
				services.CompareHashAndPassword = func(hash, plain string) error {
					return errors.New("wrong password")
				}
			},
			assert: func(t *testing.T, jwt string, err error) {
				require.Error(t, err)
				require.Zero(t, jwt)
			},
			teardown: func() {
				services.CompareHashAndPassword = compareFunc
			},
		},
	}

	loginPayload := models.LoginPayload{
		Username: "ryanpujo",
		Password: "okeokeryy5yiyykyky8yie",
	}
	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			jwt, err := credService.Login(context.Background(), &loginPayload)

			v.assert(t, jwt, err)

			v.teardown()
		})
	}
}
