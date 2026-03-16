package middleware

import (
	"enactus/internal/httpx"
	"enactus/internal/models"
	"net/http"
)

func RoleMiddleware(next httpx.AppHandler, roles ...models.Role) httpx.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		claimsValue := r.Context().Value("claims")
		if claimsValue == nil {
			return httpx.Unauthorized("claims missing")
		}

		claims, ok := claimsValue.(map[string]interface{})
		if !ok {
			return httpx.Unauthorized("claims missing")
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

		return next(w, r)
	}
}
