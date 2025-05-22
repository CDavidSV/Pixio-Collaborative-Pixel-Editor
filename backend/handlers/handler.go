package handlers

import (
	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/services"
	"github.com/CDavidSV/Pixio/websocket"
)

type Handler struct {
	queries   *data.Queries
	services  *services.Services
	websocket *websocket.Hub
}

func NewHandler(queries *data.Queries, services *services.Services, wsHub *websocket.Hub) *Handler {
	return &Handler{
		queries:   queries,
		services:  services,
		websocket: wsHub,
	}
}
