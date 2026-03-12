package middleware

import (
	"context"
	"enactus/internal/auth"
	"enactus/internal/httpx"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtSecret *auth.JWTSecret, next httpx.AppHandler) httpx.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		tokenStr := r.Header.Get("Authorization")

		if tokenStr == "" {
			return httpx.Unauthorized("empty token")
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(tokenStr, bearerPrefix) {
			return httpx.Unauthorized("invalid token format")
		}

		tokenStr = strings.TrimPrefix(tokenStr, bearerPrefix)

		token, err := jwtSecret.ValidateToken(tokenStr)
		if err != nil {
			return httpx.ValidationError(err)
		}

		ctx := context.WithValue(r.Context(), "id", token.Claims)
		return next(w, r.WithContext(ctx))
	}
}
