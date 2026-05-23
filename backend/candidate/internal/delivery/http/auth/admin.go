package auth

import (
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/helpers"
)

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := FromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if user.Role != "admin" {
			helpers.RespondError(w, http.StatusForbidden, "you are not an admin")
			return
		}
		next.ServeHTTP(w, r)
	})
}
