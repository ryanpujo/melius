package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
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
	jsonStr, _ = json.Marshal(country)
)

func TestSaveCountry(t *testing.T) {
	failed := models.Country{
		Name: "",
	}
	invalidJson, _ := json.Marshal(failed)
	tableTest := map[string]struct {
		noHeader bool
		token      string
		json       []byte
		arrange    func()
		assert     func(t *testing.T, actualCode int, res utilities.Response)
	}{
		"success": {
			token: "fhtht",
			noHeader: false,
			json:  jsonStr,
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
			json:  jsonStr,
			token: "dgg",
			noHeader: false,
			arrange: func() {
				asm.On("SaveCountry", mock.Anything, country).Return(0, errors.New("failed")).Once()
				jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
			},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusBadRequest, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "failed to record the country", res.Message)
			},
		},
		"validation error": {
			token: "dfgg",
			json:  invalidJson,
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
			json:    jsonStr,
			noHeader: false,
			arrange: func() {},
			assert: func(t *testing.T, actualCode int, res utilities.Response) {
				require.Equal(t, http.StatusUnauthorized, actualCode)
				require.NotZero(t, res)
				require.Zero(t, res.ID)
				require.Equal(t, "authentication failed", res.Message)
				require.Equal(t, "bearer token is empty", res.Err)
			},
		},
		"token verification failed": {
			json:  jsonStr,
			token: "ddgr",
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
			json: jsonStr,
			noHeader: true,
			arrange: func() {},
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
