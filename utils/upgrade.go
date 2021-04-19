package utils

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var u = upGrader() // use default options
func upGrader() websocket.Upgrader {
	return websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
}

// WebsocketConn return web socket connection
func WebsocketConn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return u.Upgrade(w, r, nil)

}
