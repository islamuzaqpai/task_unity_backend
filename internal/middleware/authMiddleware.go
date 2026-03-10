package middleware

import (
	"context"
	"enactus/internal/auth"
	"enactus/internal/httpx"
	"fmt"
	"net/http"
)

func AuthMiddleware(jwt auth.JWTSecret, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		fmt.Println(tokenStr)
		if tokenStr == "" {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		token, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "id", token.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
