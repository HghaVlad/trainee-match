package auth

import "net/http"

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := FromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if user.Role != "admin" {
			http.Error(w, "your are not admin", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
