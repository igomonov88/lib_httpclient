package lib_httpclient

import (
	"net/http"

	"github.com/eapache/go-resiliency/breaker"
)

// CircuitBreakerTransportDecorator should be last in decorators call chain
// so even if breaker is open we will still have metrics/trace/etc if configured.
func CircuitBreakerTransportDecorator(cfg CircuitBreakerConfig) TransportDecoratorFunc {
	return func(rt http.RoundTripper) http.RoundTripper {

		br := breaker.New(cfg.ErrThreshold, cfg.SuccessThreshold, cfg.Timeout)

		return TransportFunc(func(r *http.Request) (resp *http.Response, err error) {
			err = br.Run(func() error {
				resp, err = rt.RoundTrip(r)
				return err
			})

			return resp, err
		})
	}
}
