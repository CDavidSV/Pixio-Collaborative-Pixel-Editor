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

func encodeMessage(t string, m proto.Message) ([]byte, error) {
	payloadBytes, err := proto.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	message := &msg.WSMessage{
		Type:    t,
		Payload: payloadBytes,
	}
	b, err := proto.Marshal(message)
	return b, err
}

func sendError(client *WSClient, errMsg string) {
	wsError := &msg.WSError{
		Error: errMsg,
	}

	msgBytes, err := encodeMessage("error", wsError)
	if err != nil {
		slog.Error("Failed to encode error message", "Error", err.Error())
		return
	}

	client.send <- msgBytes
}

func sendMessage(client *WSClient, msgType string, msg proto.Message) {
	msgBytes, err := encodeMessage(msgType, msg)
	if err != nil {
		slog.Error("Failed to encode message", "Error", err.Error())
		return
	}

	client.send <- msgBytes
}
