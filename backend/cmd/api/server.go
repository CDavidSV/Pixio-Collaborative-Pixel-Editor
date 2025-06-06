package api

import (
	"net/http"
	"time"

	"github.com/CDavidSV/Pixio/config"
	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/handlers"
	"github.com/CDavidSV/Pixio/middlewares"
	"github.com/CDavidSV/Pixio/services"
	"github.com/CDavidSV/Pixio/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	addr     string
	queries  *data.Queries
	services *services.Services
}

func NewServer(addr string, pool *pgxpool.Pool) *Server {
	queries := data.NewQueries(pool)          // Data layer
	services := services.NewServices(queries) // Business logic layer

	return &Server{
		addr:     addr,
		queries:  queries,
		services: services,
	}
}

func (s *Server) Start() error {
	r := chi.NewRouter()

	// Create the websocket hub
	// This hub will be used to manage websocket connections and handle messages
	wsHub := websocket.NewWebsocketHub(s.queries, s.services)

	appMiddleware := middlewares.NewMiddleware(s.queries, s.services)
	handlers := handlers.NewHandler(s.queries, s.services, wsHub)

	// Middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(cors.Handler(config.CorsConfig))
	r.Use(appMiddleware.CommonHeaders)

	// Mount routes
	r.Mount("/api/v1", s.loadRoutes(handlers, appMiddleware))

	server := &http.Server{
		Addr:         s.addr,
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
