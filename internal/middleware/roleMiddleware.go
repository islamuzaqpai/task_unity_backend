package middleware

import (
	"enactus/internal/httpx"
	"enactus/internal/models"
	"net/http"
)

func RoleMiddleware(next http.Handler, roles ...models.Role) httpx.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		claimsValue := r.Context().Value("claims")
		if claimsValue == "" {
			return httpx.Unauthorized("claims missing")
		}

		claims, ok := claimsValue.(map[string]interface{})
		if !ok {
			return httpx.Unauthorized("invalid claims type")
		}

		roleValue, ok := claims["role"].(string)
		if !ok {
			return httpx.Unauthorized("role missing")
		}

		allowed := false
		for _, role := range roles {
			if roleValue == role.Name {
				allowed = true
				break
			}
		}

		if !allowed {
			return httpx.Unauthorized("insufficient permissions")
		}

		next.ServeHTTP(w, r)
		return nil
	}
}
