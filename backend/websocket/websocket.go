package websocket

import (
	"log/slog"
	"net/http"
	"sync"

	"slices"

	"github.com/CDavidSV/Pixio/config"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:       originChecker,
	EnableCompression: true,
} // use default options

type Hub struct {
	conns     map[string]*Client
	connMutex sync.RWMutex
}

func NewWebsocketHub() *Hub {
	return &Hub{
		conns: make(map[string]*Client),
	}
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

	err = h.waitForAuth(conn)

}

func (h *Hub) waitForAuth(conn *websocket.Conn) error {
	return nil
}

func originChecker(r *http.Request) bool {
	if len(config.AllowedDomains) == 0 {
		return true
	}

	return slices.Contains(config.AllowedDomains, r.Header.Get("Origin"))
}
