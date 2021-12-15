package lib_httpclient

import (
	"net/http"
)

// TransportFunc is a function type that implements the http.RoundTripper interface.
type TransportFunc func(*http.Request) (*http.Response, error)

// RoundTrip supports http.RoundTripper interface
func (f TransportFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// TransportDecoratorFunc wraps a http.RoundTripper with extra behaviour.
type TransportDecoratorFunc func(http.RoundTripper) http.RoundTripper

// DecorateTransport decorates a http.RoundTripper c with all the given Decorators, in order.
func DecorateTransport(rt http.RoundTripper, ds ...TransportDecoratorFunc) http.RoundTripper {
	result := rt
	for _, decorate := range ds {
		result = decorate(result)
	}
	return result
}

// DecorateClientTransport decorates internal Transport of http.Client
func DecorateClientTransport(c *http.Client, ds ...TransportDecoratorFunc) *http.Client {
	if c == nil {
		c = http.DefaultClient
	}

	transport := c.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	for _, decorate := range ds {
		transport = decorate(transport)
	}

	resultClient := *c
	resultClient.Transport = transport

	return &resultClient
}

func DecorateClientTransportByDefault(c *http.Client) *http.Client {
	return DecorateClientTransport(c,
		TracerTransportDecorator(),
	)
}
