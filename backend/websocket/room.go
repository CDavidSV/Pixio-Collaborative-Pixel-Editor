package websocket

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/websocket/msg"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/proto"
)

type LoadStatus uint16

const (
	NotLoaded LoadStatus = iota
	Loading
	Loaded
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
	loadStatus  LoadStatus
}

type ClientWithPerms struct {
	WSClient *WSClient
	Perms    *types.UserAccess
}

func (r *Room) SetClient(client *WSClient, accessRules *types.UserAccess) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.deleteTimer != nil {
		r.deleteTimer.Stop()
		r.deleteTimer = nil
	}

	r.Clients[client.ID] = &ClientWithPerms{
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

func (r *Room) loadCanvasData(data []byte) {
	r.loadStatus = Loading
	defer func() { r.loadStatus = Loaded }()

	pixelData, err := r.hub.services.CanvasService.LoadCanvas(data)
	if err != nil {
		slog.Error("Failed to load canvas pixel data", "Error", err.Error())
		return
	}

	r.PixelData = pixelData
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

func (h *Hub) joinRoom(client *WSClient, payload []byte) {
	joinRoom := &msg.JoinRoom{}
	if err := proto.Unmarshal(payload, joinRoom); err != nil {
		sendError(client, msg.JoinRoomMsg, ErrUnmarshallingMsg.Error())
		return
	}

	userAccess, err := h.queries.GetUserAccess(joinRoom.CanvasId, types.CanvasObject, client.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			sendError(client, msg.JoinRoomMsg, ErrMissingPermissions.Error())
			return
		}

		sendError(client, msg.JoinRoomMsg, ErrFetchingUserAccess.Error())
		return
	}

	canvas, err := h.queries.GetCanvas(joinRoom.CanvasId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			sendError(client, msg.JoinRoomMsg, ErrCanvasNotFound.Error())
		} else {
			sendError(client, msg.JoinRoomMsg, ErrFetchingCanvas.Error())
		}
		return
	}

	// Check if the room already exists
	h.roomMutex.Lock()
	room, exists := h.rooms[canvas.ID]
	if !exists {
		room = &Room{
			CanvasID:   canvas.ID,
			Clients:    make(map[string]*ClientWithPerms),
			Width:      canvas.Width,
			Height:     canvas.Height,
			hub:        h,
			loadStatus: NotLoaded,
		}
		h.rooms[canvas.ID] = room
	}

	if room.loadStatus == NotLoaded {
		go room.loadCanvasData(canvas.PixelData)
	}
	h.roomMutex.Unlock()

	room.SetClient(client, &userAccess)

	sendMessage(client, msg.JoinRoomMsg, &msg.JoinRoomSuccess{
		CanvasId: canvas.ID,
		UserId:   client.ID,
		ConnId:   client.connID,
	})
}

func (h *Hub) leaveRoom(client *WSClient, payload []byte) {

}

func (h *Hub) updateCursorPosition(client *WSClient, payload []byte) {
	mousePos := &msg.MousePosition{}
	err := proto.Unmarshal(payload, mousePos)
	if err != nil {
		sendError(client, msg.MousePosUpdateMsg, ErrUnmarshallingMsg.Error())
		return
	}

	room := client.GetRoom(mousePos.RoomId)
	if room == nil {
		sendError(client, msg.MousePosUpdateMsg, ErrRoomNotFound.Error())
		return
	}

	mousePosSend := &msg.MousePositionUpdate{
		UserId: client.ID,
		X:      mousePos.X,
		Y:      mousePos.Y,
	}

	message, err := encodeMessage(msg.MousePosUpdateMsg, mousePosSend)
	if err != nil {
		sendError(client, msg.MousePosUpdateMsg, ErrMarshallingMsg.Error())
		return
	}

	broadcastMessage(client, room, message)
}
