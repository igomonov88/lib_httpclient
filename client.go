package lib_httpclient

import "net/http"

// Client is an interface supported by http.Client.
// Use this interface whenever you need to abstract from http.Client implementation and to use decorated client.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientFunc is a function type that implements the Client interface.
type ClientFunc func(*http.Request) (*http.Response, error)

// Do supports Client interface
func (f ClientFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

// ClientDecoratorFunc wraps a Client with extra behaviour.
type ClientDecoratorFunc func(Client) Client

// DecorateClient decorates a Client c with all the given Decorators, in order.
func DecorateClient(c Client, ds ...ClientDecoratorFunc) Client {
	result := c
	for _, decorate := range ds {
		result = decorate(result)
	}
	return result
}

func DecorateClientByDefault(c Client, serviceName string) Client {
	return DecorateClient(c,
		TracerClientDecorator(serviceName),
	)
}
