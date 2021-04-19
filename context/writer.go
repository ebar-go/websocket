package context

import "github.com/gorilla/websocket"

type Writer struct {
	conn *websocket.Conn
}

func NewWriter(conn *websocket.Conn) Writer {
	return Writer{conn: conn}
}

func (w Writer) Write(p []byte) (n int, err error) {
	if err := w.conn.WriteMessage(websocket.TextMessage, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

