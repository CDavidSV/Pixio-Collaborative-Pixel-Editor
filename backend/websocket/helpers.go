package websocket

import (
	"log/slog"

	"github.com/CDavidSV/Pixio/websocket/msg"
	"google.golang.org/protobuf/proto"
)

func decodeMessage(t string, b []byte, m proto.Message) error {
	message := &msg.WSMessage{}
	err := proto.Unmarshal(b, message)
	if err != nil {
		return err
	}

	if message.Type != t {
		return ErrInvalidMsgType
	}

	return proto.Unmarshal(message.Payload, m)
}

func encodeMessage(t msg.WSMessageType, m proto.Message) ([]byte, error) {
	payloadBytes, err := proto.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	message := &msg.WSMessage{
		Type:    string(t),
		Payload: payloadBytes,
	}
	b, err := proto.Marshal(message)
	return b, err
}

func sendError(client *WSClient, msgType msg.WSMessageType, errMsg string) {
	wsError := &msg.WSError{
		Error: errMsg,
	}

	msgBytes, err := encodeMessage(msgType, wsError)
	if err != nil {
		slog.Error("Failed to encode error message", "Error", err.Error())
		return
	}

	client.send <- msgBytes
}

func sendMessage(client *WSClient, msgType msg.WSMessageType, msg proto.Message) {
	msgBytes, err := encodeMessage(msgType, msg)
	if err != nil {
		slog.Error("Failed to encode message", "Error", err.Error())
		return
	}

	client.send <- msgBytes
}

func broadcastMessage(sender *WSClient, room *Room, msg []byte) {
	room.mu.RLock()
	defer room.mu.RUnlock()
	for _, c := range room.Clients {
		// Don't send the message to the sender
		if c.WSClient.ID == sender.ID {
			continue
		}

		c.WSClient.send <- msg
	}
}
