package websocket

import (
	"log"
	"net/http"
)

// simpleServer implement of Server
type simpleServer struct {
	server
	// socket连接
	connections map[string]Connection
	// conn注册chan
	register    chan Connection
	// conn注销chan
	unregister  chan Connection

}
// HandleRequest implement of Server
func (srv *simpleServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// 获取socket连接
	conn, err := newConnection(w, r)
	if err != nil {
		// do something..
		return
	}

	srv.registerConn(conn)
}

// registerConn 注册连接
func (srv *simpleServer) registerConn(conn Connection) {
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


// Broadcast implement of Server
func (srv *simpleServer) Broadcast(response Response, ignores ...string) {
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
func (srv *simpleServer) Close(conn Connection)  {
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
func (srv *simpleServer) Start() {
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

// SampleServer 返回Server的实例
func SampleServer() Server {
	return &simpleServer{
		server: base(),
		connections: make(map[string]Connection),
		register:    make(chan Connection),
		unregister:  make(chan Connection),
	}
}