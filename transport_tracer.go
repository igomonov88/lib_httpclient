package lib_httpclient

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func TracerTransportDecorator() TransportDecoratorFunc {
	return func(rt http.RoundTripper) http.RoundTripper {
		return TransportFunc(func(r *http.Request) (*http.Response, error) {
			span, ctx := opentracing.StartSpanFromContext(r.Context(), "http_request")
			defer span.Finish()

			span.SetTag("http.host", r.URL.Host)
			span.SetTag("http.url", r.URL.String())
			span.SetTag("http.method", r.Method)

			span.SetTag("action", ActionFromRequest(r))

			r = r.WithContext(ctx)

			resp, err := rt.RoundTrip(r)
			if err != nil {
				span.SetTag("error", true)

				span.LogFields(
					log.String("event", "error"),
					log.String("message", err.Error()),
				)
			}

			if resp != nil {
				span.SetTag("http.status_code", resp.StatusCode)
			}

			return resp, err
		})
	}
}
