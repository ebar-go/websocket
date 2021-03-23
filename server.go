/**
 * @Author: Hongker
 * @Description:
 * @File:  server
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

import (
	"log"
	"net/http"
)

// Server websocket服务
type Server interface {
	// 处理请求
	HandleRequest(w http.ResponseWriter, r *http.Request)
	// 连接时触发
	HandleConnect(func (conn Connection))
	// 断开连接时触发
	HandleDisconnect(func (conn Connection))
	// 映射路由
	Route(uri string, handler func(ctx Context))
	// 关闭连接
	Close(conn Connection)
	// 广播
	Broadcast(response Response, ignoreConnections ...Connection)
	// 启动服务
	Start()
}

// server implement of Server
type server struct {
	engine *Engine
	connections map[string]Connection
	register    chan Connection
	unregister  chan Connection
}
// HandleRequest implement of Server
func (srv *server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// 获取socket连接
	conn, err := newConnection(w, r)
	if err != nil {
		// do something..
		return
	}

	srv.registerConn(conn)
}

// registerConn 注册连接
func (srv *server) registerConn(conn Connection) {
	// 开启一个协程，异步监听socket的发送
	go func() {
		// 连接断开后自动close，释放资源
		defer srv.Close(conn)
		// 使用engine监听connection
		srv.engine.listen(conn)
	}()
	// 通过channel传递connection,防止并发
	srv.register <- conn
}
// HandleConnect implement of Server
func (srv *server) HandleConnect(f func(conn Connection)) {
	panic("implement me")
}
// HandleDisconnect implement of Server
func (srv *server) HandleDisconnect(f func(conn Connection)) {
	panic("implement me")
}
// Route implement of Server
func (srv *server) Route(uri string, handler func(ctx Context)) {
	srv.engine.route(uri, handler)
}

// Broadcast implement of Server
func (srv *server) Broadcast(response Response, ignoreConnections ...Connection) {
	for id, conn := range srv.connections {
		// 跳过指定连接
		var skip bool
		for _, ignore := range ignoreConnections {
			if ignore.ID() == id {
				skip = true
				break
			}
		}
		if !skip {
			if err := conn.write(response.Byte()); err != nil {
				log.Println("write to [%s]: %v", id, err)
			}
		}
	}
}


// Close implement of Server
func (srv *server) Close(conn Connection)  {
	// 关闭socket
	conn.close()
	// 注销conn
	srv.unregister <- conn
}

// Start implement of Server
func (srv *server) Start() {
	// 设置默认的404路由
	if srv.engine.noRoute == nil {
		srv.engine.NoRoute(notFoundHandler)
	}
	go func() {
		for {
			select {
			case conn := <-srv.register: // 注册connection
				srv.connections[conn.ID()] = conn
			case conn := <-srv.unregister: // 注销connection
				delete(srv.connections, conn.ID())
			}
		}
	}()
}

var _default = New()

// Default 使用默认实例
func Default() Server {
	return _default
}
// New 返回Server的实例
func New() Server {
	return &server{
		engine: &Engine{routers: map[string]Handler{}},
		connections: make(map[string]Connection),
		register:    make(chan Connection),
		unregister:  make(chan Connection),
	}
}