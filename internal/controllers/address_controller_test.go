package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/utilities"
	"github.com/ryanpujo/melius/proof"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type addressServiceMock struct {
	mock.Mock
}

func (asm *addressServiceMock) SaveCountry(ctx context.Context, country models.Country) (uint, error) {
	args := asm.Called(ctx, country)
	return uint(args.Int(0)), args.Error(1)
}

func (asm *addressServiceMock) SaveState(ctx context.Context, state models.State, countryID uint) (uint, error) {
	args := asm.Called(ctx, state, countryID)
	return uint(args.Int(0)), args.Error(1)
}

func (asm *addressServiceMock) SaveCity(ctx context.Context, city models.City, stateID uint) (uint, error) {
	args := asm.Called(ctx, city, stateID)
	return uint(args.Int(0)), args.Error(1)
}

func (asm *addressServiceMock) SaveAddress(ctx context.Context, address models.Address, cityID uint) (uint, error) {
	args := asm.Called(ctx, address, cityID)
	return uint(args.Int(0)), args.Error(1)
}

var (
	country = models.Country{
		Name: "Indonesia",
	}
	state = models.State{
		Name:    "Jakarta",
		Country: country,
	}
	city = models.City{
		Name:  "Jakarta Timur",
		State: state,
	}
	address = models.Address{
		AddressLine: "jl. mayjen sutoyo kel. cawang kec kramat jati rt.007/011",
		PostalCode:  "12630",
		IsMain:      true,
		City:        city,
	}
	token = &jwt.Token{
		Claims: jwt.MapClaims{
			"username": "test",
		},
	}
)

func TestSaveCountry(t *testing.T) {
	jsonStr, _ := json.Marshal(country)
	failed := models.Country{
		Name: "",
	}
	invalidJson, _ := json.Marshal(failed)
	tableTest := map[string]struct {
		noHeader bool
		token    string
		json     []byte
		arrange  func()
		assert   func(t *testing.T, actualCode int, res utilities.Response)
	}{
		"success": {
			token:    "fhtht",
			noHeader: false,
			json:     jsonStr,
			arrange: func() {
				asm.On("SaveCountry", mock.Anything, country).Return(1, nil).Once()
				jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusCreated, actualCode)
				require.NotZero(t, res)
				require.Equal(t, uint(1), res.ID)
			},
		},
		"failed": {
			json:     jsonStr,
			token:    "dgg",
			noHeader: false,
			arrange: func() {
				asm.On("SaveCountry", mock.Anything, country).Return(0, errors.New("failed")).Once()
				jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusBadRequest, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "Failed to record the country", res.Message)
			},
		},
		"validation error": {
			token:    "dfgg",
			json:     invalidJson,
			noHeader: false,
			arrange: func() {
				jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusBadRequest, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "Validation Error", res.Message)
			},
		},
		"empty bearer token": {
			json:     jsonStr,
			noHeader: false,
			arrange:  func() {},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusUnauthorized, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "authentication failed", res.Message)
				require.Equal(t, "bearer token is empty", res.Err)
			},
		},
		"token verification failed": {
			json:     jsonStr,
			token:    "ddgr",
			noHeader: false,
			arrange: func() {
				jwtm.On("VerifyToken", mock.Anything).Return((*jwt.Token)(nil), errors.New("verification is failed")).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusUnauthorized, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "authentication failed", res.Message)
				require.Equal(t, "verification is failed", res.Err)
			},
		},
		"no authorization header": {
			json:     jsonStr,
			noHeader: true,
			arrange:  func() {},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusUnauthorized, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "authentication failed", res.Message)
				require.Equal(t, "Token is required", res.Err)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			res, code, err := proof.NewHttpProof(http.MethodPost,
				"/auth/country",
				proof.WithJSON(v.json),
				proof.WithJWTToken(v.token),
				proof.WithNoAuthorizationHeader(v.noHeader),
			).RunTest(handler)
			require.NoError(t, err)

			v.assert(t, code, res)
		})
	}
}

func TestSaveState(t *testing.T) {
	jsonStr, _ := json.Marshal(state)
	failedState := models.State{
		Name: "jakarta",
	}
	invalidJson, _ := json.Marshal(&failedState)
	tableTest := map[string]struct {
		pathVar  uint
		noHeader bool
		token    string
		json     []byte
		arrange  func()
		assert   func(t *testing.T, actualCode int, res utilities.Response)
	}{
		"success": {
			pathVar:  1,
			token:    "ghfththth",
			json:     jsonStr,
			noHeader: false,
			arrange: func() {
				asm.On("SaveState", mock.Anything, state, uint(1)).Return(1, nil).Once()
				jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusCreated, actualCode)
				require.NotZero(t, res)
				require.Equal(t, uint(1), res.ID)
				jwtm.AssertCalled(t, "VerifyToken", "ghfththth")
			},
		},
		"failed": {
			pathVar:  1,
			token:    "dgdrjrngk",
			json:     jsonStr,
			noHeader: false,
			arrange: func() {
				asm.On("SaveState", mock.Anything, state, uint(1)).Return(0, errors.New("failed")).Once()
				jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusBadRequest, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "Failed to record the state", res.Message)
			},
		},
		"validation error": {
			json:     invalidJson,
			token:    "dgrgrg",
			pathVar:  1,
			noHeader: false,
			arrange: func() {
				jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusBadRequest, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "Validation Error", res.Message)
			},
		},
		"no path variable": {
			pathVar: 0,
			json: jsonStr,
			token: "dgrgrg",
			noHeader: false,
			arrange: func() {
				jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusBadRequest, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "No country associated with this state", res.Message)
			},
		},
		"empty bearer token": {
			token: "",
			pathVar: 1,
			noHeader: false,
			json: jsonStr,
			arrange: func() {},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusUnauthorized, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "authentication failed", res.Message)
				require.Equal(t, "bearer token is empty", res.Err)
			},
		},
		"no authorization header": {
			pathVar: 1,
			noHeader: true,
			token: "",
			json: jsonStr,
			arrange: func() {},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusUnauthorized, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "authentication failed", res.Message)
				require.Equal(t, "Token is required", res.Err)
			},
		},
		"failed to verify token": {
			pathVar: 1,
			noHeader: false,
			token: "dgrgrg",
			json: jsonStr,
			arrange: func() {
				jwtm.On("VerifyToken", mock.Anything).Return((*jwt.Token)(nil), errors.New("failed to verify")).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusUnauthorized, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "authentication failed", res.Message)
				require.Equal(t, "failed to verify", res.Err)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			res, code, err := proof.NewHttpProof(http.MethodPost,
				fmt.Sprintf("/auth/state/%d", v.pathVar),
				proof.WithJSON(v.json),
				proof.WithJWTToken(v.token),
				proof.WithNoAuthorizationHeader(v.noHeader),
			).RunTest(handler)
			require.NoError(t, err)

			v.assert(t, code, res)
		})
	}
}
