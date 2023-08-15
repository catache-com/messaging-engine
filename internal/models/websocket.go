package models

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Config
const (
	readBufferBytesSize  = 4096
	writeBufferBytesSize = 4096
)

// makeUpgrader: make a websocket upgrader based on specified buffer sizes.
// Utility function for UpgradeHTTPToWS function.
func makeUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  readBufferBytesSize,
		WriteBufferSize: writeBufferBytesSize,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
}

// UpgradeHTTPToWS upgrades the HTTP server connection to the WebSocket protocol.
func UpgradeHTTPToWS(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := makeUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, err
}
