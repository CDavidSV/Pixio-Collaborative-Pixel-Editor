package websocket

import (
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"slices"

	"github.com/CDavidSV/Pixio/config"
	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/services"
	"github.com/CDavidSV/Pixio/websocket/msg"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:       originChecker,
	EnableCompression: true,
} // use default options

type HandlerFunc func(client *WSClient, payload []byte)

type Hub struct {
	conns     map[string]*WSClient
	connMutex sync.RWMutex
	queries   *data.Queries
	services  *services.Services
	handlers  map[string]HandlerFunc
	rooms     map[string]*Room
	roomMutex sync.RWMutex
}

type WSClient struct {
	ID   string
	send chan []byte
	conn *websocket.Conn
}

func NewWebsocketHub(queries *data.Queries, services *services.Services) *Hub {
	hub := &Hub{
		conns:    make(map[string]*WSClient),
		queries:  queries,
		services: services,
		handlers: make(map[string]HandlerFunc),
		rooms:    make(map[string]*Room),
	}
	hub.registerHandlers()

	return hub
}

func (h *Hub) registerHandlers() {
}

func (h *Hub) WSHanlder(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Error upgrading connection: ", "error", err.Error())
		return
	}

	if err := conn.SetCompressionLevel(4); err != nil {
		slog.Error("Error setting compression level: ", "error", err.Error())
		return
	}

	client := &WSClient{
		send: make(chan []byte),
		conn: conn,
	}

	userID, err := h.waitForAuth(conn)
	if err != nil {
		slog.Error("Error authenticating user", "Error", err.Error())

		if errors.Is(err, ErrFailedAuth) {
			sendError(client, err.Error())
		} else if errors.Is(err, ErrInvalidMsgType) {
			sendError(client, err.Error())
		} else {
			sendError(client, "UNEXPECTED_SERVER_ERROR")
		}

		conn.Close()
		return
	}

	client.ID = userID

	h.addConnection(client)
	defer h.removeConnection(client)

	go h.writePump(client)
	h.readPump(client)
}

func (h *Hub) writePump(c *WSClient) {
	for data := range c.send {
		c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		err := c.conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			slog.Error("Error writing message to client: ", "error", err.Error())
			break
		}
	}
}

func (h *Hub) readPump(c *WSClient) {
	for {
		msgType, msgData, err := c.conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message: ", "error", err.Error())
			return
		}

		if msgType != websocket.BinaryMessage {
			continue
		}

		message := &msg.WSMessage{}
		err = proto.Unmarshal(msgData, message)
		if err != nil {
			sendError(c, "ERROR_DECODING_MESSAGE")
		}
	}
}

func (h *Hub) waitForAuth(conn *websocket.Conn) (string, error) {
	msgType, msgData, err := conn.ReadMessage()
	if err != nil {
		return "", err
	}

	if msgType != websocket.BinaryMessage {
		return "", ErrInvalidMsgType
	}

	authMsg := &msg.Auth{}
	err = decodeMessage("auth", msgData, authMsg)
	if err != nil {
		return "", err
	}

	// validate access token
	userID, ok := h.services.AuthService.ValidateAccessToken(authMsg.Token)
	if !ok {
		return userID, ErrFailedAuth
	}

	return userID, nil
}

func (h *Hub) addConnection(client *WSClient) {
	slog.Info("new connection established", "ip", client.conn.RemoteAddr().String())

	h.connMutex.Lock()
	defer h.connMutex.Unlock()

	h.conns[client.ID] = client
}

func (h *Hub) removeConnection(client *WSClient) {
	slog.Info("connection closed", "ip", client.conn.RemoteAddr().String())

	h.roomMutex.RLock()
	for _, room := range h.rooms {
		room.RemoveClient(client.ID)
	}
	h.roomMutex.RUnlock()

	close(client.send)
	client.conn.Close()

	h.connMutex.Lock()
	delete(h.conns, client.ID)
	h.connMutex.Unlock()
}

func (h *Hub) getClient(userID string) (*WSClient, bool) {
	h.connMutex.RLock()
	defer h.connMutex.RUnlock()

	client, ok := h.conns[userID]
	return client, ok
}

func originChecker(r *http.Request) bool {
	if len(config.AllowedDomains) == 0 {
		return true
	}

	return slices.Contains(config.AllowedDomains, r.Header.Get("Origin"))
}
