package main

import (
	"github.com/CDavidSV/Pixio/cmd/api/config"
	"github.com/CDavidSV/Pixio/cmd/api/handlers"
	"github.com/CDavidSV/Pixio/cmd/api/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func loadRoutes() *chi.Mux {
	r := chi.NewRouter()

	handlers := handlers.NewAppHandlers()
	middleware := middleware.NewAppMiddleware()

	// Middleware
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.CleanPath)
	r.Use(cors.Handler(config.CorsConfig))
	r.Use(middleware.CommonHeaders)

	// Routes
	r.Route("/auth", func(r chi.Router) {
		r.Use(chiMiddleware.AllowContentType("application/x-www-form-urlencoded"))

		r.Post("/signup", handlers.Signup)
		r.Post("/login", handlers.Login)
		r.Post("/token", handlers.Token)
		r.Post("/logout", handlers.Logout)
	})

	return r
}
