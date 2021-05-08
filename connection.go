/**
 * @Author: Hongker
 * @Description:
 * @File:  connection
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:06
 */

package websocket

import (
	"github.com/ebar-go/websocket/context"
	"github.com/ebar-go/websocket/utils"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// Connection websocket连接
type Connection interface {
	// 唯一标识
	ID() string
	// 给客户端发送数据
	Write(msg []byte) error
	// 关闭连接
	close() error
	// 获取上下文
	context() (Context, error)
	// 获取连接的文件描述符
	fd() int
}

// connection 自定义websocket连接
type connection struct {
	// 连接唯一ID，可以用于业务逻辑设计
	id string
	// socket connection
	sockConn *websocket.Conn
	// socket连接的文件标识符
	sockFD int
}

// ID implement of Connection
func (conn *connection) ID() string {
	return conn.id
}

// write implement of Connection
func (conn *connection) Write(msg []byte) error {
	return conn.sockConn.WriteMessage(websocket.TextMessage, msg)
}

// close implement of Connection
func (conn *connection) close() error {
	return conn.sockConn.Close()
}

// context implement of Connection
func (conn *connection) context() (Context, error) {
	ctx := context.NewContext(conn.sockConn)
	if err := ctx.Read(); err != nil {
		return nil, err
	}
	return ctx, nil
}

// fd 获取文件标识符
func (conn *connection) fd() int {
	return conn.sockFD
}

// newConnection return initialized websocket Connection
func newConnection(w http.ResponseWriter, r *http.Request) (*connection, error) {
	socketConn, err := utils.WebsocketConn(w, r)
	if err != nil {
		return nil, err
	}

	conn := &connection{
		id:       uuid.NewV4().String(),
		sockConn: socketConn,
		sockFD:   utils.SocketFD(socketConn.UnderlyingConn()),
	}
	return conn, nil
}
