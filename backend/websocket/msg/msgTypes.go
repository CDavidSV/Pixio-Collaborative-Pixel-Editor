package msg

type WSMessageType string

var (
	ErrorMsg          WSMessageType = "error"
	AuthMsg           WSMessageType = "auth"
	MousePosUpdateMsg WSMessageType = "mouse_position_update"
)
