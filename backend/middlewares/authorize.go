package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/CDavidSV/Pixio/utils"
)

func (m *Middleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		s := strings.Split(authorization, " ")
		if len(s) != 2 || s[0] != "Bearer" || s[1] == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID, ok := m.services.AuthService.ValidateAccessToken(s[1])
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
