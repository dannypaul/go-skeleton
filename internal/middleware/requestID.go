package middleware

import (
	"context"
	"net/http"

	"github.com/dannypaul/go-skeleton/internal/kit/http/header"

	"github.com/google/uuid"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(header.RequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set(header.RequestID, requestID)

		ctx = context.WithValue(ctx, "requestId", uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
