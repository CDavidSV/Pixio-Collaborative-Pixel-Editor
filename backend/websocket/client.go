package websocket

import (
	"sync"

	"github.com/CDavidSV/Pixio/utils"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	ID          string
	connID      string
	send        chan []byte
	conn        *websocket.Conn
	joinedRooms map[string]*Room
	mu          sync.RWMutex
}

func NewClient(userID string, conn *websocket.Conn) *WSClient {
	return &WSClient{
		ID:          userID,
		connID:      utils.GenerateID(),
		conn:        conn,
		send:        make(chan []byte),
		joinedRooms: make(map[string]*Room),
		mu:          sync.RWMutex{},
	}
}

func (c *WSClient) AddRoom(room *Room) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.joinedRooms[room.CanvasID] = room
}

func (c *WSClient) RemoveRoom(roomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.joinedRooms != nil {
		delete(c.joinedRooms, roomID)
	}
}

func (c *WSClient) GetRoom(roomID string) *Room {
	c.mu.RLock()
	defer c.mu.RUnlock()

	room, ok := c.joinedRooms[roomID]
	if !ok {
		return nil
	}
	return room
}
