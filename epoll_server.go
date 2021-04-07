/**
 * @Author: Hongker
 * @Description:
 * @File:  epollServer
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

import (
	"log"
	"net/http"
)


// epollServer implement of Server
type epollServer struct {
	server

	epoller *epoll
}
// HandleRequest implement of Server
func (srv *epollServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// 获取socket连接
	conn, err := newConnection(w, r)
	if err != nil {
		// do something..
		return
	}

	srv.registerConn(conn)
}

// registerConn 注册连接
func (srv *epollServer) registerConn(conn Connection) {
	if err := srv.epoller.Add(conn); err != nil {
		log.Printf("Failed to add connection")
		conn.close()
		return
	}
	// 注册回调
	if srv.connectCallback != nil {
		srv.connectCallback(conn)
	}
}

// Broadcast implement of Server
func (srv *epollServer) Broadcast(response Response, ignores ...string) {
	for _, conn := range srv.epoller.connections {
		// 跳过指定连接
		var skip bool
		for _, ignore := range ignores {
			if ignore == conn.ID() {
				skip = true
				break
			}
		}
		if !skip {
			if err := conn.write(response.Byte()); err != nil {
				log.Printf("write to [%s]: %v", conn.ID(), err)
			}
		}
	}
}


// Close implement of Server
func (srv *epollServer) Close(conn Connection)  {
	if err := srv.epoller.Remove(conn); err != nil {
		log.Printf("Failed to remove %v", err)
	}
	// 关闭socket
	conn.close()
	// 注销回调
	if srv.disconnectCallback != nil {
		srv.disconnectCallback(conn)
	}
}

// Start implement of Server
func (srv *epollServer) Start() {

	// epoll模式
	go func() {
		for {
			connections, err := srv.epoller.Wait()
			if err != nil {
				log.Printf("Failed to epoll wait %v", err)
				continue
			}
			for _, conn := range connections {
				ctx, err := conn.context()
				if err != nil {
					srv.Close(conn)
					continue
				}
				srv.engine.handle(ctx)
			}
		}
	}()
}
// EpollServer 通过epoll模式实现的websocket服务
func EpollServer() Server {
	epoller, err := MkEpoll()
	if err != nil {
		log.Fatalf("create epoll:%v\n", err)
	}
	return &epollServer{
		server: base(),
		epoller: epoller,
	}
}