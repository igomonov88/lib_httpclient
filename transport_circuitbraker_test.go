package lib_httpclient

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransportWithCircuitBreaker(t *testing.T) {
	mockTransport := &testTransport{numFailures: 3, t: t, expectedBody: []byte("body")}

	cfg := CircuitBreakerConfig{2, 1, time.Minute}

	transport := DecorateTransport(mockTransport, CircuitBreakerTransportDecorator(cfg))
	client := &http.Client{Transport: transport}

	// call 3 times with 3 errs in resp
	for i := 1; i <= 3; i++ {
		req := buildReq(t, mockTransport.expectedBody)
		_, err := client.Do(req)
		require.Error(t, err)
	}

	// call 4 time, client will return ok result
	// but cb should not call it at all
	req := buildReq(t, mockTransport.expectedBody)
	_, _ = client.Do(req)
	//require.Equal(t, "Get http://127.0.0.1/: circuit breaker is open", err.Error())
}

type testTransport struct {
	numFailures  int
	expectedBody []byte
	t            *testing.T
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	assert.Equal(t.t, "GET", req.Method)
	assert.Equal(t.t, "http://127.0.0.1/", req.URL.String())
	if t.expectedBody != nil {
		defer req.Body.Close()
		b, err := ioutil.ReadAll(req.Body)
		assert.Nil(t.t, err)
		assert.Equal(t.t, string(t.expectedBody), string(b))
	} else {
		assert.Nil(t.t, req.Body)
	}

	if t.numFailures <= 0 {
		return &http.Response{}, nil
	}

	t.numFailures--

	return nil, errors.New("unexpected error")
}
