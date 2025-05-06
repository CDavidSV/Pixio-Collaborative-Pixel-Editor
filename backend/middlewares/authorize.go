package middlewares

import (
	"net/http"
	"strings"
)

func (m *Middleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		s := strings.Split(authorization, " ")
		if len(s) != 2 || s[0] != "Bearer" || s[1] == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !m.services.AuthService.ValidAccessToken(s[1]) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
