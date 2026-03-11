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
			return httpx.BadRequest("invalid token")
		}

		const bearerPrefix = "Bearer "
		tokenStr = strings.TrimSpace(tokenStr)
		if len(tokenStr) > len(bearerPrefix) && tokenStr[:len(bearerPrefix)] == bearerPrefix {
			tokenStr = tokenStr[len(bearerPrefix):]
		}

		token, err := jwtSecret.ValidateToken(tokenStr)
		if err != nil {
			return httpx.ValidationError(err)
		}

		ctx := context.WithValue(r.Context(), "id", token.Claims)
		return next(w, r.WithContext(ctx))
	}
}
