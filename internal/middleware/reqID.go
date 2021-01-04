package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "requestId", uuid.New().String())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
