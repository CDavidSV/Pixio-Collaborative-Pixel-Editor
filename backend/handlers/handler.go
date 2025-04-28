package handlers

import (
	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/services"
)

type Handler struct {
	queries  *data.Queries
	services *services.Services
}

func NewHandler(queries *data.Queries, services *services.Services) *Handler {
	return &Handler{
		queries:  queries,
		services: services,
	}
}
