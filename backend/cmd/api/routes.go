package api

import (
	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/handlers"
	"github.com/CDavidSV/Pixio/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) loadRoutes(queries *data.Queries, services *services.Services) *chi.Mux {
	handlers := handlers.NewHandler(queries, services)

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
