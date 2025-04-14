package middleware

import "net/http"

type AppMiddleware struct {
}

func NewAppMiddleware() *AppMiddleware {
	return &AppMiddleware{}
}

func (m *AppMiddleware) CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}
