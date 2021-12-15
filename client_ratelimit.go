package lib_httpclient

import (
	"net/http"
	"time"
)

func RateLimitClientDecorator(interval time.Duration) ClientDecoratorFunc {
	return func(c Client) Client {
		if interval <= 0 {
			return c
		}

		ticker := time.NewTicker(interval)

		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			<-ticker.C
			return c.Do(r)
		})
	}
}
