package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dannypaul/go-skeleton/internal/config"
	"github.com/dannypaul/go-skeleton/internal/iam"
	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware decodes the token in Authorization header and packs it into context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf, err := config.Get()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var claims iam.Claims
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ctx := context.WithValue(r.Context(), iam.CtxClaimsKey, claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		var authToken string
		splitHeader := strings.Split(authHeader, "Bearer ")
		if len(splitHeader) > 1 {
			authToken = splitHeader[1]
		}
		token, err := jwt.ParseWithClaims(authToken, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.JwtSecret), nil
		})

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		_, ok := token.Claims.(*iam.Claims)
		if !ok || !token.Valid || claims.UserId == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), iam.CtxClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
