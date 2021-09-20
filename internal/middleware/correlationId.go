package middleware

import (
	"context"
	"net/http"

	"github.com/dannypaul/go-skeleton/internal/kit/http/header"

	"github.com/google/uuid"
)

func CorrelationId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		correlationID := r.Header.Get(header.CorrelationId)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		w.Header().Set(header.CorrelationId, correlationID)

		ctx = context.WithValue(ctx, "correlationId", uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
