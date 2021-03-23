/**
 * @Author: Hongker
 * @Description:
 * @File:  connection
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:06
 */

package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type connection struct {
	id string
	// socket connection
	sockConn *websocket.Conn

}

func (conn *connection) ID() string {
	return conn.id
}

// writeMessage 发送数据
func (conn *connection) write(msg []byte) error{
	return conn.sockConn.WriteMessage(websocket.TextMessage, msg)
}

func (conn *connection) close(unregister chan <- Connection) {
	_ = conn.sockConn.Close()
	unregister <- conn
}

// Listen listen connection
func (c *connection) listen(engine *Engine) {

	for {
		ctx, err := c.context()
		if err != nil {
			break
		}

		engine.Run(ctx)
	}
}

func (c *connection) context() (Context, error) {
	_, message, err := c.sockConn.ReadMessage()
	if err != nil {
		return nil, err
	}
	request := new(MessageRequest)
	if err := json.Unmarshal(message, request); err != nil {
		// 参数错误
		//return nil, err
	}
	ctx := &context{request: request, conn: c}
	return ctx, nil
}


var u = upGrader() // use default options
func upGrader() websocket.Upgrader {
	return websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
}
// WebsocketConn return web socket connection
func newConnection(w http.ResponseWriter, r *http.Request) (Connection, error) {
	respHeader := http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}}
	conn, err := u.Upgrade(w, r, respHeader)
	if err != nil {
		return nil, err
	}

	return &connection{
		id:   uuid.NewV4().String(),
		sockConn: conn,
	}, nil
}
