package middlewares

import (
	"net/http"

	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/services"
)

type Middleware struct {
	queries  *data.Queries
	services *services.Services
}

func NewMiddleware(queries *data.Queries, services *services.Services) *Middleware {
	return &Middleware{
		queries:  queries,
		services: services,
	}
}

func (m *Middleware) CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}
