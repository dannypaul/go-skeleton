package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const RequestIDHeader = "X-Request-Id"

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, "requestId", uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
