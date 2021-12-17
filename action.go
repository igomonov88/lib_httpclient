package lib_httpclient

import (
	"context"
	"net/http"
)

func RequestWithAction(req *http.Request, action string) *http.Request {
	return req.WithContext(ContextWithAction(req.Context(), action))
}

func ContextWithAction(ctx context.Context, action string) context.Context {
	return context.WithValue(ctx, contextKeyAction, action)
}

func ActionFromRequest(req *http.Request) string {
	action, _ := req.Context().Value(contextKeyAction).(string)
	return action
}
