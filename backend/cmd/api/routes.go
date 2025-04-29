package api

import (
	"github.com/CDavidSV/Pixio/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) loadRoutes() *chi.Mux {
	handlers := handlers.NewHandler(s.queries, s.services)

	r := chi.NewRouter()

	// Authentication routes
	r.Route("/auth", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/x-www-form-urlencoded"))

		r.Post("/signup", handlers.Signup)
		r.Post("/login", handlers.Login)
		r.Post("/token", handlers.Token)
		r.Post("/logout", handlers.Logout)
	})

	return r
}
