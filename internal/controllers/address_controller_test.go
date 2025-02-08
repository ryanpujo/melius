package controllers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/utilities"
	"github.com/ryanpujo/melius/proof"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type addressServiceMock struct {
	mock.Mock
}

func (asm *addressServiceMock) SaveCountry(ctx context.Context, country *models.CountryPayload) (*models.Country, error) {
	args := asm.Called(ctx, country)
	return args.Get(0).(*models.Country), args.Error(1)
}

func (asm *addressServiceMock) SaveState(ctx context.Context, state *models.StatePayload, countryID uint) (*models.State, error) {
	args := asm.Called(ctx, state, countryID)
	return args.Get(0).(*models.State), args.Error(1)
}

func (asm *addressServiceMock) SaveCity(ctx context.Context, city *models.CityPayload, stateID uint) (*models.City, error) {
	args := asm.Called(ctx, city, stateID)
	return args.Get(0).(*models.City), args.Error(1)
}

func (asm *addressServiceMock) SaveAddress(ctx context.Context, address *models.AddressPayload) (*models.Address, error) {
	args := asm.Called(ctx, address)
	return args.Get(0).(*models.Address), args.Error(1)
}

func (asm *addressServiceMock) GetCountries(ctx context.Context) ([]*models.Country, error) {
	args := asm.Called(ctx)
	return args.Get(0).([]*models.Country), args.Error(1)
}

func (asm *addressServiceMock) GetCityByID(ctx context.Context, cityID uint) (*models.City, error) {
	args := asm.Called(ctx, cityID)
	return args.Get(0).(*models.City), args.Error(1)
}

func (asm *addressServiceMock) CreateUserAddress(ctx context.Context, address *models.AddressPayload, userID uint) (*models.Address, error) {
	args := asm.Called(ctx, address, userID)
	return args.Get(0).(*models.Address), args.Error(1)
}

var (
	countryPayload = &models.CountryPayload{
		Name: "Indonesia",
	}
	country = &models.Country{
		ID:   1,
		Name: "Indonesia",
	}

	statePayload = models.StatePayload{
		Name: "Jakarta",
	}
	state = &models.State{
		ID:   1,
		Name: "Jakarta",
	}

	cityPayload = models.CityPayload{
		Name: "Jakarta Timur",
	}
	city = &models.City{
		ID:   1,
		Name: "Jakarta Timur",
	}

	addressPayload = models.AddressPayload{
		AddressLine: "jl. mayjen sutoyo kel. cawang kec kramat jati rt.007/011",
		PostalCode:  "12630",
		IsMain:      true,
		CityID:      uint(1),
	}
	address = &models.Address{
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
	authURI          = "/auth"
	emptyBearerToken = testCase{
		noHeader: false,
		arrange:  func() {},
		assert: func(t *testing.T, actualCode int, res utilities.Response) {
			require.Equal(t, http.StatusUnauthorized, actualCode)
			require.NotZero(t, res)
			require.Equal(t, "authentication failed", res.Message)
			require.Equal(t, "bearer token is empty", res.Err)
		},
	}
	noAuthorizationHeader = testCase{
		noHeader: true,
		arrange:  func() {},
		assert: func(t *testing.T, actualCode int, res utilities.Response) {
			require.Equal(t, http.StatusUnauthorized, actualCode)
			require.NotZero(t, res)
			require.Equal(t, "authentication failed", res.Message)
			require.Equal(t, "Token is required", res.Err)
		},
	}
	tokenVerificationFailed = testCase{
		token:    "ddgr",
		noHeader: false,
		arrange: func() {
			jwtm.On("VerifyToken", mock.Anything).Return((*jwt.Token)(nil), errors.New("verification is failed")).Once()
		},
		assert: func(t *testing.T, actualCode int, res utilities.Response) {
			require.Equal(t, http.StatusUnauthorized, actualCode)
			require.NotZero(t, res)
			require.Equal(t, "authentication failed", res.Message)
			require.Equal(t, "verification is failed", res.Err)
		},
	}
	authTestCases = testCases{
		"empty bearer token":        emptyBearerToken,
		"no authorization header":   noAuthorizationHeader,
		"token verification failed": tokenVerificationFailed,
	}
)

type testCase struct {
	pathVar  uint
	noHeader bool
	token    string
	json     []byte
	arrange  func()
	assert   func(t *testing.T, actualCode int, res utilities.Response)
}
type testCases map[string]testCase

func withAuthTestCases(specificCases testCases) testCases {
	allCases := make(testCases)

	for k, v := range authTestCases {
		allCases[k] = v
	}

	for k, v := range specificCases {
		allCases[k] = v
	}

	return allCases
}

func TestSaveCountry(t *testing.T) {
	jsonStr, _ := json.Marshal(countryPayload)
	failed := models.CountryPayload{
		Name: "",
	}
	invalidJson, _ := json.Marshal(failed)
	tableTest := withAuthTestCases(
		testCases{
			"success": {
				token:    "fhtht",
				noHeader: false,
				json:     jsonStr,
				arrange: func() {
					asm.On("SaveCountry", mock.Anything, countryPayload).Return(country, nil).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusCreated, actualCode)
					require.NotZero(t, res)
					require.Equal(t, country, res.Country)
				},
			},
			"failed": {
				json:     jsonStr,
				token:    "dgg",
				noHeader: false,
				arrange: func() {
					asm.On("SaveCountry", mock.Anything, countryPayload).
						Return((*models.Country)(nil), errors.New("failed")).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusInternalServerError, actualCode)
					require.NotZero(t, res)
					require.Zero(t, res.Country)
					require.Equal(t, "An internal server error occurred. Please try again later.", res.Message)
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
					require.Zero(t, res.Country)
					require.Equal(t, "There was a problem with your request. Please double-check your input and try again.", res.Message)
				},
			},
		},
	)
	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			res, code, err := proof.NewHttpProof(http.MethodPost,
				fmt.Sprintf("%s%s", authURI, "/country"),
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
	jsonStr, _ := json.Marshal(statePayload)
	failedState := models.StatePayload{
		Name: "",
	}
	invalidJson, _ := json.Marshal(&failedState)
	tableTest := withAuthTestCases(
		testCases{
			"success": {
				pathVar:  1,
				token:    "ghfththth",
				json:     jsonStr,
				noHeader: false,
				arrange: func() {
					asm.On("SaveState", mock.Anything, &statePayload, uint(1)).Return(state, nil).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusCreated, actualCode)
					require.NotZero(t, res)
					require.Equal(t, state, res.State)
					jwtm.AssertCalled(t, "VerifyToken", "ghfththth")
				},
			},
			"failed": {
				pathVar:  1,
				token:    "dgdrjrngk",
				json:     jsonStr,
				noHeader: false,
				arrange: func() {
					asm.On("SaveState", mock.Anything, &statePayload, uint(1)).Return((*models.State)(nil), errors.New("failed")).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusInternalServerError, actualCode)
					require.NotZero(t, res)
					require.Zero(t, res.State)
					require.Equal(t, "An internal server error occurred. Please try again later.", res.Message)
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
					require.Zero(t, res.State)
					require.Equal(t, "There was a problem with your request. Please double-check your input and try again.", res.Message)
				},
			},
			"no path variable": {
				pathVar:  0,
				json:     jsonStr,
				token:    "dgrgrg",
				noHeader: false,
				arrange: func() {
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusBadRequest, actualCode)
					require.NotZero(t, res)
					require.Zero(t, res.State)
					require.Equal(t, "There was a problem with your request. Please double-check your input and try again.", res.Message)
				},
			},
		},
	)

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			res, code, err := proof.NewHttpProof(http.MethodPost,
				fmt.Sprintf("%s%s", authURI, fmt.Sprintf("/state/%d", v.pathVar)),
				proof.WithJSON(v.json),
				proof.WithJWTToken(v.token),
				proof.WithNoAuthorizationHeader(v.noHeader),
			).RunTest(handler)
			require.NoError(t, err)

			v.assert(t, code, res)
		})
	}
}

func TestSaveCity(t *testing.T) {
	jsonStr, _ := json.Marshal(cityPayload)
	failedCity := models.CityPayload{
		Name: "",
	}
	invalidJson, _ := json.Marshal(&failedCity)
	tableTest := withAuthTestCases(
		testCases{
			"success": {
				pathVar:  1,
				token:    "ghfththth",
				json:     jsonStr,
				noHeader: false,
				arrange: func() {
					asm.On("SaveCity", mock.Anything, &cityPayload, uint(1)).Return(city, nil).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusCreated, actualCode)
					require.NotZero(t, res)
					require.Equal(t, city, res.City)
					jwtm.AssertCalled(t, "VerifyToken", "ghfththth")
				},
			},
			"failed": {
				pathVar:  1,
				token:    "dgdrjrngk",
				json:     jsonStr,
				noHeader: false,
				arrange: func() {
					asm.On("SaveCity", mock.Anything, &cityPayload, uint(1)).Return((*models.City)(nil), errors.New("failed")).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusInternalServerError, actualCode)
					require.NotZero(t, res)
					require.Zero(t, res.City)
					require.Equal(t, "An internal server error occurred. Please try again later.", res.Message)
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
					require.Zero(t, res.City)
					require.Equal(t, "There was a problem with your request. Please double-check your input and try again.", res.Message)
				},
			},
			"no path variable": {
				pathVar:  0,
				json:     jsonStr,
				token:    "dgrgrg",
				noHeader: false,
				arrange: func() {
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusBadRequest, actualCode)
					require.NotZero(t, res)
					require.Zero(t, res.City)
					require.Equal(t, "There was a problem with your request. Please double-check your input and try again.", res.Message)
				},
			},
		},
	)

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			res, code, err := proof.NewHttpProof(http.MethodPost,
				fmt.Sprintf("%s%s", authURI, fmt.Sprintf("/city/%d", v.pathVar)),
				proof.WithJSON(v.json),
				proof.WithJWTToken(v.token),
				proof.WithNoAuthorizationHeader(v.noHeader),
			).RunTest(handler)
			require.NoError(t, err)

			v.assert(t, code, res)
		})
	}
}

func TestSaveAddress(t *testing.T) {
	jsonStr, _ := json.Marshal(addressPayload)
	failedAddress := models.AddressPayload{
		AddressLine: "dgrgrgrg",
	}
	invalidJson, _ := json.Marshal(&failedAddress)
	tableTest := withAuthTestCases(
		testCases{
			"success": {
				token:    "ghfththth",
				json:     jsonStr,
				noHeader: false,
				arrange: func() {
					asm.On("SaveAddress", mock.Anything, &addressPayload).Return(address, nil).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusCreated, actualCode)
					require.NotZero(t, res)
					require.Equal(t, address, res.Address)
					jwtm.AssertCalled(t, "VerifyToken", "ghfththth")
				},
			},
			"failed": {
				token:    "dgdrjrngk",
				json:     jsonStr,
				noHeader: false,
				arrange: func() {
					asm.On("SaveAddress", mock.Anything, &addressPayload).Return((*models.Address)(nil), errors.New("failed")).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusInternalServerError, actualCode)
					require.NotZero(t, res)
					require.Zero(t, res.Address)
					require.Equal(t, "An internal server error occurred. Please try again later.", res.Message)
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
					require.Zero(t, res.Address)
					require.Equal(t, "There was a problem with your request. Please double-check your input and try again.", res.Message)
				},
			},
		},
	)

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			res, code, err := proof.NewHttpProof(http.MethodPost,
				fmt.Sprintf("%s%s", authURI, "/address"),
				proof.WithJSON(v.json),
				proof.WithJWTToken(v.token),
				proof.WithNoAuthorizationHeader(v.noHeader),
			).RunTest(handler)
			require.NoError(t, err)

			v.assert(t, code, res)
		})
	}
}

func TestGetCountries(t *testing.T) {
	tableTest := withAuthTestCases(
		testCases{
			"success": {
				token:    "sdfdfsf",
				noHeader: false,
				arrange: func() {
					asm.On("GetCountries", mock.Anything).Return([]*models.Country{country}, nil).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusOK, actualCode)
					require.Equal(t, res.Countries[0], country)
					require.Len(t, res.Countries, 1)
				},
			},
			"failed": {
				token:    "ggrgr",
				noHeader: false,
				arrange: func() {
					asm.On("GetCountries", mock.Anything).Return(([]*models.Country)(nil), errors.New("failed")).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusInternalServerError, actualCode)
					require.Equal(t, "An internal server error occurred. Please try again later.", res.Message)
					require.Zero(t, res.Countries)
				},
			},
		},
	)

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			res, code, err := proof.NewHttpProof(http.MethodGet,
				fmt.Sprintf("%s%s", authURI, "/country"),
				proof.WithJWTToken(v.token),
				proof.WithNoAuthorizationHeader(v.noHeader),
			).RunTest(handler)
			require.NoError(t, err)

			v.assert(t, code, res)
		})
	}
}

func TestGetCityByID(t *testing.T) {
	tableTest := withAuthTestCases(
		testCases{
			"success": {
				pathVar:  1,
				token:    "grrgrr",
				noHeader: false,
				arrange: func() {
					asm.On("GetCityByID", mock.Anything, uint(1)).Return(city, nil).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusOK, actualCode)
					require.NotNil(t, res)
					require.Equal(t, city, res.City)
				},
			},
			"failed": {
				pathVar:  1,
				token:    "sgfef",
				noHeader: false,
				arrange: func() {
					asm.On("GetCityByID", mock.Anything, uint(1)).
						Return((*models.City)(nil), utilities.HandleError(&pgconn.PgError{Code: utilities.PGNotNullViolation})).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusBadRequest, actualCode)
					require.Zero(t, res.City)
					require.Equal(t, "A required field is missing. Please check your input.", res.Message)
				},
			},
			"err no rows error": {
				pathVar:  1,
				token:    "sgfef",
				noHeader: false,
				arrange: func() {
					asm.On("GetCityByID", mock.Anything, uint(1)).
						Return((*models.City)(nil), utilities.HandleError(sql.ErrNoRows)).Once()
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusNotFound, actualCode)
					require.Zero(t, res.City)
					require.Equal(t, "We're sorry, but we couldn't find any data with that information. Please check your input and try again.", res.Message)
				},
			},
			"no path variable": {
				token:    "grgrg",
				noHeader: false,
				arrange: func() {
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusBadRequest, actualCode)
					require.Zero(t, res.City)
					require.Equal(t, "There was a problem with your request. Please double-check your input and try again.", res.Message)
				},
			},
		},
	)

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			res, code, err := proof.NewHttpProof(http.MethodGet,
				fmt.Sprintf("%s%s", authURI, fmt.Sprintf("/city/%d", v.pathVar)),
				proof.WithJWTToken(v.token),
				proof.WithNoAuthorizationHeader(v.noHeader),
			).RunTest(handler)
			require.NoError(t, err)

			v.assert(t, code, res)
		})
	}
}

func TestCreateUserAddress(t *testing.T) {
	jsonStr, _ := json.Marshal(addressPayload)
	failedAddress := models.AddressPayload{
		AddressLine: "dgrgrgrg",
	}
	invalidJson, _ := json.Marshal(&failedAddress)
	tableTest := withAuthTestCases(
		testCases{
			"success": {
				token:    "sff",
				pathVar:  1,
				noHeader: false,
				json:     jsonStr,
				arrange: func() {
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
					asm.On("CreateUserAddress", mock.Anything, &addressPayload, uint(1)).
						Return(address, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusCreated, actualCode)
					require.Equal(t, address, res.Address)
				},
			},
			"failed": {
				token:    "sff",
				pathVar:  1,
				noHeader: false,
				json:     jsonStr,
				arrange: func() {
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
					asm.On("CreateUserAddress", mock.Anything, &addressPayload, uint(1)).
						Return((*models.Address)(nil), errors.New("failed")).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusInternalServerError, actualCode)
					require.Equal(t, "An internal server error occurred. Please try again later.", res.Message)
				},
			},
			"uri validation error": {
				token:    "sff",
				noHeader: false,
				json:     jsonStr,
				arrange: func() {
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusBadRequest, actualCode)
					require.Equal(t, "There was a problem with your request. Please double-check your input and try again.", res.Message)
				},
			},
			"json validation error": {
				token:    "sff",
				pathVar:  1,
				noHeader: false,
				json:     invalidJson,
				arrange: func() {
					jwtm.On("VerifyToken", mock.Anything).Return(token, nil).Once()
				},
				assert: func(t *testing.T, actualCode int, res utilities.Response) {
					require.Equal(t, http.StatusBadRequest, actualCode)
					require.Equal(t, "There was a problem with your request. Please double-check your input and try again.", res.Message)
				},
			},
		},
	)

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			res, code, err := proof.NewHttpProof(http.MethodPost,
				fmt.Sprintf("%s/user/%d/add-address", authURI, v.pathVar),
				proof.WithJSON(v.json),
				proof.WithJWTToken(v.token),
				proof.WithNoAuthorizationHeader(v.noHeader),
			).RunTest(handler)
			require.NoError(t, err)

			v.assert(t, code, res)
		})
	}
}
