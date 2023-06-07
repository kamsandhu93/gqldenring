package middleware

import (
	"context"
	"net/http"
	"time"
)

func WithLogging(logFn func(context.Context, string, ...any), h http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		logFn(r.Context(), "Incoming request %s %s", uri, method)

		h.ServeHTTP(rw, r) // serve the original request

		duration := time.Since(start)
		logFn(r.Context(), "Request complete%s", duration)
	}
	return http.HandlerFunc(fn)
}
