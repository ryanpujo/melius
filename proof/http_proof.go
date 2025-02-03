package proof

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/ryanpujo/melius/config"
	"github.com/ryanpujo/melius/internal/utilities"
)

type httpProof struct {
	JSON                  []byte
	httpMethod            string
	uri                   string
	jwtToken              string
	noAuthorizationHeader bool
}

type proofFunc func(*httpProof)

func WithNoAuthorizationHeader(header bool) proofFunc {
	return func(hp *httpProof) {
		hp.noAuthorizationHeader = header
	}
}

func WithJSON(jsonStr []byte) proofFunc {
	return func(hp *httpProof) {
		hp.JSON = jsonStr
	}
}

func WithJWTToken(token string) proofFunc {
	return func(hp *httpProof) {
		hp.jwtToken = token
	}
}

func NewHttpProof(httpMethod string, uri string, opts ...proofFunc) *httpProof {
	proof := &httpProof{
		httpMethod:            httpMethod,
		uri:                   uri,
		noAuthorizationHeader: false,
	}

	for _, opt := range opts {
		opt(proof)
	}

	return proof
}

func (hp *httpProof) RunTest(handler http.Handler) (utilities.Response, int, error) {
	var reader io.Reader

	if hp.JSON != nil {
		reader = bytes.NewReader(hp.JSON)
	}

	fullURI := fmt.Sprintf("%s%s", config.Config().BaseRoute, hp.uri)

	req, err := http.NewRequest(hp.httpMethod, fullURI, reader)
	if err != nil {
		return utilities.Response{}, 0, err
	}

	if !hp.noAuthorizationHeader {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", hp.jwtToken))
	}

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var res utilities.Response

	err = json.NewDecoder(rec.Body).Decode(&res)
	if err != nil {
		return utilities.Response{}, 0, err
	}

	return res, rec.Code, nil
}
