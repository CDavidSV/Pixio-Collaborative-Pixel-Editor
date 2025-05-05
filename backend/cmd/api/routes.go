package api

import (
	"github.com/CDavidSV/Pixio/handlers"
	"github.com/CDavidSV/Pixio/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) loadRoutes(handlers *handlers.Handler, appMiddleware *middlewares.Middleware) *chi.Mux {
	r := chi.NewRouter()

	// Authentication routes
	r.Route("/auth", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/x-www-form-urlencoded"))

		r.Post("/signup", handlers.SignupPost)
		r.Post("/login", handlers.LoginPost)
		r.Post("/token", handlers.TokenPost)
		r.Post("/logout", handlers.LogoutPost)
	})

	return r
}
