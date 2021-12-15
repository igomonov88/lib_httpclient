package lib_httpclient

import (
	"net/http"
	"time"

	"github.com/eapache/go-resiliency/breaker"
)

type CircuitBreakerConfig struct {
	// Num of successive errors to breaker became opened
	ErrThreshold int
	// Num of successful call to breaker became closed from half-opened
	SuccessThreshold int
	// Timeout to breaker became half-opened after it became opened
	Timeout time.Duration
}

// CircuitBreakerClientDecorator should be last in decorators call chain
// so even if breaker is open we will still have metrics/trace/etc if configured.
func CircuitBreakerClientDecorator(cfg CircuitBreakerConfig) ClientDecoratorFunc {
	return func(c Client) Client {

		br := breaker.New(cfg.ErrThreshold, cfg.SuccessThreshold, cfg.Timeout)

		return ClientFunc(func(r *http.Request) (resp *http.Response, err error) {
			err = br.Run(func() error {
				resp, err = c.Do(r)
				return err
			})

			return resp, err
		})
	}
}
