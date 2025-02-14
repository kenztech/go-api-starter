package middlewares

import (
	"net/http"

	"github.com/kenztech/go-api-starter/utils"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userEmail, userRole, err := utils.VerifyRefreshToken(cookie.Value)
		if err != nil {
			utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := utils.SetUserDataInContext(r.Context(), userEmail, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly is a middleware that ensures the user has the 'admin' role
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, ok := utils.GetUserDataFromContext(r.Context())
		if !ok || userData.Role != "admin" {
			utils.SendError(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
