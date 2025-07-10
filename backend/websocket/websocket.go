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
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:       originChecker,
	EnableCompression: true,
} // use default options

type HandlerFunc func(client *WSClient, payload []byte)

type Hub struct {
	conns     map[string]map[string]*WSClient // userID -> connID -> client
	connMutex sync.RWMutex
	queries   *data.Queries
	services  *services.Services
	handlers  map[string]HandlerFunc
	rooms     map[string]*Room
	roomMutex sync.RWMutex
}

func NewWebsocketHub(queries *data.Queries, services *services.Services) *Hub {
	hub := &Hub{
		conns:    make(map[string]map[string]*WSClient),
		queries:  queries,
		services: services,
		handlers: make(map[string]HandlerFunc),
		rooms:    make(map[string]*Room),
	}
	hub.registerHandlers()

	return hub
}

func (h *Hub) registerHandlers() {
	h.handlers[string(msg.MousePosUpdateMsg)] = h.updateCursorPosition
	h.handlers[string(msg.JoinRoomMsg)] = h.joinRoom
	h.handlers[string(msg.LeaveRoomMsg)] = h.leaveRoom
}

func (h *Hub) WSHanlder(w http.ResponseWriter, r *http.Request) {
	// get the user id from the request params
	userID := chi.URLParam(r, "id")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Error upgrading connection: ", "error", err.Error())
		return
	}

	if err := conn.SetCompressionLevel(4); err != nil {
		slog.Error("Error setting compression level: ", "error", err.Error())
		return
	}

	client := NewClient(userID, conn)

	err = h.waitForAuth(conn, userID)
	if err != nil {
		slog.Error("Error authenticating user", "Error", err.Error())

		if errors.Is(err, ErrFailedAuth) {
			sendError(client, msg.ErrorMsg, err.Error())
		} else if errors.Is(err, ErrInvalidMsgType) {
			sendError(client, msg.ErrorMsg, err.Error())
		} else {
			sendError(client, msg.ErrorMsg, ErrUnexpected.Error())
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
			sendError(c, msg.ErrorMsg, ErrDecodingMsg.Error())
		}

		h.executeHandler(c, message)
	}
}

func (h *Hub) executeHandler(client *WSClient, message *msg.WSMessage) {
	handler, ok := h.handlers[message.Type]
	if !ok {
		sendError(client, msg.ErrorMsg, ErrUnsupportedMsgType.Error())
		return
	}

	handler(client, message.Payload)
}

func (h *Hub) waitForAuth(conn *websocket.Conn, providedUserID string) error {
	msgType, msgData, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	if msgType != websocket.BinaryMessage {
		return ErrInvalidMsgType
	}

	authMsg := &msg.Auth{}
	err = decodeMessage("auth", msgData, authMsg)
	if err != nil {
		return err
	}

	// validate access token
	userID, ok := h.services.AuthService.ValidateAccessToken(authMsg.Token)
	if !ok {
		return ErrFailedAuth
	}

	if userID != providedUserID {
		return ErrFailedAuth
	}

	return nil
}

func (h *Hub) addConnection(client *WSClient) {
	slog.Info("new connection established", "ip", client.conn.RemoteAddr().String())

	h.connMutex.Lock()
	defer h.connMutex.Unlock()

	if h.conns[client.ID] == nil {
		h.conns[client.ID] = make(map[string]*WSClient)
	}

	h.conns[client.ID][client.connID] = client
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

func (h *Hub) getClient(userID, connID string) (*WSClient, bool) {
	h.connMutex.RLock()
	defer h.connMutex.RUnlock()

	client, ok := h.conns[userID][connID]
	return client, ok
}

func originChecker(r *http.Request) bool {
	if len(config.AllowedDomains) == 0 {
		return true
	}

	return slices.Contains(config.AllowedDomains, r.Header.Get("Origin"))
}
