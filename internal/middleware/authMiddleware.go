package middleware

import (
	"context"
	"enactus/internal/auth"
	"enactus/internal/httpx"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return httpx.Unauthorized("invalid token claims")
		}

		userId := int(claims["user_id"].(float64))
		role, ok := claims["role"].(string)
		if !ok {
			return httpx.Unauthorized("role missing in token")
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		ctx = context.WithValue(ctx, "claims", map[string]interface{}{
			"user_id": userId,
			"role":    role,
		})
		return next(w, r.WithContext(ctx))
	}
}
