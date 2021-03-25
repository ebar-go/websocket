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
	"github.com/ebar-go/websocket/utils"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

// Connection websocket连接
type Connection interface {
	// 连接的唯一ID
	ID() string
	// 给客户端发送数据
	write(msg []byte) error
	// 关闭连接
	close()
	// 获取上下文
	context() (Context, error)
	// fd
	fd() int
}

// connection 自定义websocket连接
type connection struct {
	// uuid
	id string
	// socket connection
	sockConn *websocket.Conn
}
// ID implement of Connection
func (conn *connection) ID() string {
	return conn.id
}

// write implement of Connection
func (conn *connection) write(msg []byte) error{
	return conn.sockConn.WriteMessage(websocket.TextMessage, msg)
}
// close implement of Connection
func (conn *connection) close() {
	if err := conn.sockConn.Close(); err != nil {
		log.Println("close connection: %v", err)
	}
}

// context implement of Connection
func (conn *connection) context() (Context, error) {
	_, message, err := conn.sockConn.ReadMessage()
	if err != nil {
		return nil, err
	}
	req := new(request)
	if err := json.Unmarshal(message, req); err != nil {
		// 参数错误
		log.Printf("unmarshal request: %v, source: %s \n", err, string(message))
	}
	ctx := &context{request: req, conn: conn}
	return ctx, nil
}

// fd
func (conn *connection) fd() int {
	return utils.SocketFD(conn.sockConn.UnderlyingConn())
}
// newConnection
func newConnection(w http.ResponseWriter, r *http.Request) (Connection, error) {
	conn, err := utils.WebsocketConn(w, r)
	if err != nil {
		return nil, err
	}
	return &connection{
		id:   uuid.NewV4().String(),
		sockConn: conn,
	}, nil
}
