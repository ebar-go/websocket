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
	Route(uri string, handler Handler)
	// 关闭连接
	Close(conn Connection)
	// 广播
	Broadcast(response Response, ignores ...string)
	// 启动服务
	Start()
}

// server implement of Server
type server struct {
	// 路由引擎
	engine *Engine
	// socket连接
	connections map[string]Connection
	// conn注册chan
	register    chan Connection
	// conn注销chan
	unregister  chan Connection
	// 连接回调
	connectCallback func(conn Connection)
	// 注销回调
	disconnectCallback func(conn Connection)
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
	// TODO 用epoll优化
	// 开启一个协程，异步监听socket的发送
	go func() {
		// 连接断开后自动close，释放资源
		defer srv.Close(conn)
		// 使用engine监听connection
		srv.engine.listen(conn)
	}()
	// 注册回调
	if srv.connectCallback != nil {
		srv.connectCallback(conn)
	}
	// 通过channel传递connection,防止并发
	srv.register <- conn
}
// HandleConnect implement of Server
func (srv *server) HandleConnect(callback func(conn Connection)) {
	srv.connectCallback = callback
}
// HandleDisconnect implement of Server
func (srv *server) HandleDisconnect(callback func(conn Connection)) {
	srv.disconnectCallback = callback
}
// Route implement of Server
func (srv *server) Route(uri string, handler Handler) {
	srv.engine.route(uri, handler)
}

// Broadcast implement of Server
func (srv *server) Broadcast(response Response, ignores ...string) {
	for id, conn := range srv.connections {
		// 跳过指定连接
		var skip bool
		for _, ignore := range ignores {
			if ignore == id {
				skip = true
				break
			}
		}
		if !skip {
			if err := conn.write(response.Byte()); err != nil {
				log.Printf("write to [%s]: %v", id, err)
			}
		}
	}
}


// Close implement of Server
func (srv *server) Close(conn Connection)  {
	// 关闭socket
	conn.close()
	// 注销回调
	if srv.disconnectCallback != nil {
		srv.disconnectCallback(conn)
	}
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