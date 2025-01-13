package controllers_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanpujo/melius/internal/adapter"
	"github.com/ryanpujo/melius/internal/controllers"
	"github.com/ryanpujo/melius/internal/jwttoken"
	"github.com/ryanpujo/melius/internal/route"
	"github.com/stretchr/testify/mock"
)

type jwtMock struct {
	mock.Mock
}

func (j *jwtMock) VerifyToken(token string) (*jwt.Token, error) {
	args := j.Called(token)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

var (
	csm     *CredServiceMock
	asm     *addressServiceMock
	jwtm    *jwtMock
	handler http.Handler
)

func TestMain(m *testing.M) {
	csm = new(CredServiceMock)
	asm = new(addressServiceMock)
	jwtm = new(jwtMock)
	credController := controllers.NewCredentialController(csm)
	addressController := controllers.NewAddressController(asm)
	handlerFunc := adapter.Adapter{
		CredentialController: credController,
		AddressController:    addressController,
	}

	handler = route.SetupRoutes(&handlerFunc, jwttoken.GetJWTAuth(jwtm))

	os.Exit(m.Run())
}
