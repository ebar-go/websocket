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
	"log"
	"net"
	"net/http"
	"reflect"
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


func (conn *connection) fd() int {

	return socketFD(conn.sockConn.UnderlyingConn())
}

func socketFD(conn net.Conn) int {
	//tls := reflect.TypeOf(conn.UnderlyingConn()) == reflect.TypeOf(&tls.Conn{})
	// Extract the file descriptor associated with the connection
	//connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	//if tls {
	//	tcpConn = reflect.Indirect(tcpConn.Elem())
	//}
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
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
