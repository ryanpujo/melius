package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ryanpujo/melius/internal/adapter"
	"github.com/ryanpujo/melius/internal/controllers"
	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/route"
	"github.com/ryanpujo/melius/internal/utilities"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type CredServiceMock struct {
	mock.Mock
}

func (crm *CredServiceMock) Write(ctx context.Context, payload models.UserPayload) (uint, error) {
	args := crm.Called(ctx, payload)
	return uint(args.Int(0)), args.Error(1)
}

func (csm *CredServiceMock) FindByUsername(ctx context.Context, username string) (*models.Credential, error) {
	return nil, nil
}

func (csm *CredServiceMock) Login(ctx context.Context, payload *models.LoginPayload) (string, error) {
	args := csm.Called(ctx, payload)
	return args.String(0), args.Error(1)
}

var (
	csm     *CredServiceMock
	handler http.Handler
)

func TestMain(m *testing.M) {
	csm = new(CredServiceMock)
	credController := controllers.NewCredentialController(csm)

	handlerFunc := adapter.Adapter{
		CredentialController: credController,
	}

	handler = route.SetupRoutes(&handlerFunc)

	os.Exit(m.Run())
}

func TestWrite(t *testing.T) {
	credential := models.CredentialPayload{
		Email:    "ryanpujo@gmail.com",
		Username: "ryanpujo",
		Password: "okeoke",
	}
	invalidCredential := models.CredentialPayload{
		Email:    "ryanpujogmail.com",
		Username: "ryanpujo",
		Password: "okeoke",
	}
	user := models.UserPayload{
		FirstName:         "Ryan",
		LastName:          "Pujo",
		CredentialPayload: credential,
	}
	validJson, _ := json.Marshal(user)
	invalidJson, _ := json.Marshal(invalidCredential)
	tableTest := map[string]struct {
		json    []byte
		arrange func()
		assert  func(t *testing.T, statusCode int, json utilities.RegistrationResponse)
	}{
		"success": {
			json: validJson,
			arrange: func() {
				csm.On("Write", mock.Anything, user).Return(1, nil).Once()
			},
			assert: func(t *testing.T, statusCode int, json utilities.RegistrationResponse) {
				require.Equal(t, http.StatusCreated, statusCode)
				require.Equal(t, uint(1), json.ID)
			},
		},
		"failed": {
			json: validJson,
			arrange: func() {
				csm.On("Write", mock.Anything, user).Return(0, errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, json utilities.RegistrationResponse) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Zero(t, json.ID)
				require.Equal(t, "Failed to create user", json.Message)
			},
		},
		"validation failed": {
			json:    invalidJson,
			arrange: func() {},
			assert: func(t *testing.T, statusCode int, json utilities.RegistrationResponse) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Equal(t, "Validation error", json.Message)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			req := httptest.NewRequest(http.MethodPost, "/regis", bytes.NewReader(v.json))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var jsonRes utilities.RegistrationResponse

			json.NewDecoder(res.Body).Decode(&jsonRes)

			v.assert(t, res.Code, jsonRes)
		})
	}
}

func TestLogin(t *testing.T) {
	loginPayload := models.LoginPayload{
		Username: "ryanpujo",
		Password: "okeoke",
	}
	invalidLoginPayload := models.LoginPayload{
		Username: "",
		Password: "okeoke",
	}
	jsonStrValid, _ := json.Marshal(loginPayload)
	invalidJson, _ := json.Marshal(invalidLoginPayload)
	tableTest := map[string]struct {
		json    []byte
		arrange func()
		assert  func(t *testing.T, statusCode int, json utilities.RegistrationResponse)
	}{
		"success": {
			json: jsonStrValid,
			arrange: func() {
				csm.On("Login", mock.Anything, mock.Anything).Return("token", nil).Once()
			},
			assert: func(t *testing.T, statusCode int, json utilities.RegistrationResponse) {
				require.Equal(t, http.StatusOK, statusCode)
				require.NotNil(t, json)
				require.NotZero(t, json.Token)
			},
		},
		"failed": {
			json: jsonStrValid,
			arrange: func() {
				csm.On("Login", mock.Anything, mock.Anything).Return("", errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, json utilities.RegistrationResponse) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, json)
				require.Zero(t, json.Token)
				require.Equal(t, "Login failed", json.Message)
			},
		},
		"validation failed": {
			json:    invalidJson,
			arrange: func() {},
			assert: func(t *testing.T, statusCode int, json utilities.RegistrationResponse) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Equal(t, "Validation error", json.Message)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(v.json))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var jsonRes utilities.RegistrationResponse

			json.NewDecoder(res.Body).Decode(&jsonRes)

			v.assert(t, res.Code, jsonRes)
		})
	}
}
