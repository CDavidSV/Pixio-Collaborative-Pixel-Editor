package websocket

import (
	"errors"
)

var (
	ErrFailedAuth         = errors.New("FAILED_AUTH")
	ErrInvalidMsgType     = errors.New("INVALID_MSG_TYPE")
	ErrUnmarshalling      = errors.New("UNMARSHALLING_ERROR")
	ErrMissingPermissions = errors.New("MISSING_PERMISSIONS")
	ErrCanvasNotFound     = errors.New("CANVAS_NOT_FOUND")
	ErrFetchingCanvas     = errors.New("CANNOT_FETCH_CANVAS")
	ErrLoadingCanvas      = errors.New("LOADING_CANVAS_FAILED")
	ErrClientNotConnected = errors.New("CLIENT_NOT_CONNECTED")
)
