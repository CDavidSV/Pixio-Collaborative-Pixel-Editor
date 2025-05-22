package websocket

import (
	"log/slog"
	"sync"
	"time"

	"github.com/CDavidSV/Pixio/types"
)

type Room struct {
	CanvasID    string
	Clients     map[string]*ClientWithPerms
	Width       uint16
	Height      uint16
	PixelData   []types.Pixel
	mu          sync.RWMutex
	hub         *Hub
	deleteTimer *time.Timer
}

type ClientWithPerms struct {
	WSClient *WSClient
	Perms    *types.UserAccess
}

func (r *Room) SetClient(clientID string, client *WSClient, accessRules *types.UserAccess) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.deleteTimer != nil {
		r.deleteTimer.Stop()
		r.deleteTimer = nil
	}

	r.Clients[clientID] = &ClientWithPerms{
		WSClient: client,
		Perms:    accessRules,
	}
}

func (r *Room) RemoveClient(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Clients, clientID)

	if len(r.Clients) == 0 {
		r.deleteEmtyRoomTimer(r.CanvasID)
	}
}

func (r *Room) deleteEmtyRoomTimer(roomID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.deleteTimer != nil {
		r.deleteTimer.Stop()
		r.deleteTimer = nil
	}

	r.deleteTimer = time.AfterFunc(time.Minute*3, func() {
		r.hub.roomMutex.Lock()
		defer r.hub.roomMutex.Unlock()

		room, exists := r.hub.rooms[roomID]
		if !exists {
			return
		}

		room.mu.RLock()
		defer room.mu.RUnlock()

		if len(room.Clients) == 0 {
			delete(r.hub.rooms, roomID)
		}
	})
}

func (h *Hub) JoinRoom(roomID, clientID string, userAccess types.UserAccess) error {
	canvas, err := h.queries.GetCanvas(roomID)
	if err != nil {
		return err
	}

	// Check if the room already exists
	h.roomMutex.Lock()
	room, exists := h.rooms[canvas.ID]
	if !exists {
		pixelData, err := h.services.CanvasService.LoadCanvas(canvas.PixelData)
		if err != nil {
			h.roomMutex.Unlock()
			slog.Error("Failed to load canvas pixel data", "Error", err.Error())
			return err
		}

		room = &Room{
			CanvasID:  canvas.ID,
			Clients:   make(map[string]*ClientWithPerms),
			Width:     canvas.Width,
			Height:    canvas.Height,
			PixelData: pixelData,
			hub:       h,
		}
		h.rooms[canvas.ID] = room

	}
	h.roomMutex.Unlock()

	client, ok := h.getClient(clientID)
	if !ok {
		return ErrClientNotConnected
	}

	room.SetClient(clientID, client, &userAccess)

	return nil
}

func (h *Hub) LeaveRoom(roomID, clientID string) {
	h.roomMutex.RLock()
	defer h.roomMutex.RUnlock()

	room, exists := h.rooms[roomID]
	if !exists {
		return
	}

	room.RemoveClient(clientID)
}
