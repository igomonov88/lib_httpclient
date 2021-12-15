package lib_httpclient

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
)

func TracerClientDecorator(serviceName string) ClientDecoratorFunc {
	return func(c Client) Client {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			span, ctx := opentracing.StartSpanFromContext(r.Context(), serviceName)
			defer span.Finish()

			span.SetTag("action", ActionFromRequest(r))

			r = r.WithContext(ctx)

			resp, err := c.Do(r)
			if err != nil {
				span.SetTag("error", true)
			}

			return resp, err
		})
	}
}
