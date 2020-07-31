package app

import (
	"context"
	"net/http"
	"os"

	"github.com/ComedicChimera/tempest/server/models"
	"github.com/dgrijalva/jwt-go"
)

func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path

		// allow for login path and nonce paths
		if reqPath == "/login" || reqPath == "/get-nonce" {
			next.ServeHTTP(w, r)
			return
		}

		tokHeader := r.Header.Get("Authorization")

		if tokHeader == "" {
			http.Error(w, "Missing authentication token", http.StatusForbidden)
			return
		}

		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TEMPEST_TOKEN_PWD")), nil
		})

		// invalid token
		if err != nil || !token.Valid {
			http.Error(w, "Token is not valid", http.StatusForbidden)
			return
		}

		// create token context
		ctx := context.WithValue(r.Context(), "user", tk.UserID)

		// proceed as normal
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
