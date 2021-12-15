package lib_httpclient

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testClient struct {
	numFailures  int
	expectedBody []byte
	t            *testing.T
}

func (c *testClient) Do(req *http.Request) (*http.Response, error) {
	assert.Equal(c.t, "GET", req.Method)
	assert.Equal(c.t, "http://127.0.0.1/", req.URL.String())
	if c.expectedBody != nil {
		defer req.Body.Close()
		b, err := ioutil.ReadAll(req.Body)
		assert.Nil(c.t, err)
		assert.Equal(c.t, string(c.expectedBody), string(b))
	} else {
		assert.Nil(c.t, req.Body)
	}

	if c.numFailures <= 0 {
		return &http.Response{}, nil
	}

	c.numFailures--

	return nil, errors.New("unexpected error")
}

func TestClientCircuitBreaker(t *testing.T) {
	mockClient := &testClient{numFailures: 3, t: t, expectedBody: []byte("body")}

	cfg := CircuitBreakerConfig{2, 1, time.Minute}

	client := DecorateClient(
		mockClient,
		CircuitBreakerClientDecorator(cfg),
	)

	// call 3 times with 3 errs in resp
	for i := 1; i <= 3; i++ {
		req := buildReq(t, mockClient.expectedBody)
		_, err := client.Do(req)
		require.Error(t, err)
	}

	// call 4 time, client will return ok result
	// but cb should not call it at all
	req := buildReq(t, mockClient.expectedBody)
	_, err := client.Do(req)
	require.Equal(t, breaker.ErrBreakerOpen, err)
}

func buildReq(t *testing.T, expectedBody []byte) *http.Request {
	req, err := http.NewRequest(
		"GET",
		"http://127.0.0.1/",
		bytes.NewReader(expectedBody))
	require.Nil(t, err)
	return req
}
