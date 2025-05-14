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

		r.Post("/signup", handlers.PostSignup)
		r.Post("/login", handlers.PostLogin)
		r.Post("/token", handlers.PostToken)
		r.Post("/logout", handlers.PostLogout)
	})

	// Canvas routes
	r.Route("/canvas", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(appMiddleware.Authorize)
		r.Use(appMiddleware.AuthorizeCanvasAccess)

		r.Post("/create", handlers.PostCreateCanvas)
		r.Get("/{id}", handlers.GetCanvas)
	})

	// User access routes
	r.Route("/access", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(appMiddleware.Authorize)

		r.Post("/{id}/create", handlers.PostCreateAccess)
		r.Post("/{id}/delete", handlers.PostDeleteAccess)
		r.Put("/{id}/update", handlers.PutUpdateAccess)
		r.Get("/{id}", handlers.GetAccessRules)
		r.Put("/", handlers.PutUpdateGlobalAccess)
	})

	return r
}
