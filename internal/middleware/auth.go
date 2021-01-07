package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dannypaul/go-skeleton/internal/config"
	"github.com/dannypaul/go-skeleton/internal/iam"
	"github.com/dannypaul/go-skeleton/internal/kit/http/header"

	"github.com/dgrijalva/jwt-go"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf, err := config.Get()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var claims iam.Claims
		authHeader := r.Header.Get(header.Authorization)
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
