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
	user := models.UserPayload{
		FirstName:         "Ryan",
		LastName:          "Pujo",
		CredentialPayload: credential,
	}
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, statusCode int, json utilities.RegistrationResponse)
	}{
		"success": {
			arrange: func() {
				csm.On("Write", mock.Anything, user).Return(1, nil).Once()
			},
			assert: func(t *testing.T, statusCode int, json utilities.RegistrationResponse) {
				require.Equal(t, http.StatusCreated, statusCode)
				require.Equal(t, uint(1), json.ID)
			},
		},
		"failed": {
			arrange: func() {
				csm.On("Write", mock.Anything, user).Return(0, errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, json utilities.RegistrationResponse) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Zero(t, json.ID)
				require.Equal(t, "failed to create user", json.Message)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()
			jsonReq, _ := json.Marshal(user)

			req := httptest.NewRequest(http.MethodPost, "/regis", bytes.NewReader(jsonReq))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var jsonRes utilities.RegistrationResponse

			json.NewDecoder(res.Body).Decode(&jsonRes)

			v.assert(t, res.Code, jsonRes)
		})
	}
}
