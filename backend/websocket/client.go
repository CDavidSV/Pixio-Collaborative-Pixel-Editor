package websocket

import "github.com/gorilla/websocket"

type Client struct {
	ID   string
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn, userID string) *Client {
	return &Client{
		conn: conn,
		ID:   userID,
	}
}

func (c *Client) writePump() {

}

func (c *Client) readPump() {

}
