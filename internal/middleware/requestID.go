package middleware

import (
	"context"
	"github.com/dannypaul/go-skeleton/internal/header"
	"github.com/google/uuid"
	"net/http"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(header.RequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, "requestId", uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
