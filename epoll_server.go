/**
 * @Author: Hongker
 * @Description:
 * @File:  epollServer
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

import (
	"github.com/ebar-go/websocket/epoll"
	"log"
	"net/http"
)


// epollServer implement of Server
type epollServer struct {
	server

	epoller epoll.Epoll
}
// HandleRequest implement of Server
func (srv *epollServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// 获取socket连接
	conn, err := newConnection(w, r)
	if err != nil {
		// do something..
		return
	}

	// epoll add socket fd
	if err := srv.epoller.Add(conn.fd()); err != nil {
		log.Printf("Failed to add connection")
		_ = conn.close()
		return
	}

	srv.AddConnection(conn)


}

// Close implement of Server
func (srv *epollServer) Close(conn Connection)  {
	if err := srv.epoller.Remove(conn.fd()); err != nil {
		log.Printf("Failed to remove %v", err)
	}

	srv.RemoveConnection(conn)
}

// Start implement of Server
func (srv *epollServer) Start() {
	// epoll模式
	go func() {
		for {
			// active connections
			fds, err := srv.epoller.Wait()
			if err != nil {
				log.Println("unable to get active socket connection from epoll:", err)
				continue
			}
			// handle context
			for _, fd := range fds {
				conn, exist := srv.GetConnection(fd)
				if !exist {
					continue
				}

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
	e, err := epoll.Create()
	if err != nil {
		log.Fatalf("create epoll:%v\n", err)
	}
	return &epollServer{
		server: base(),
		epoller: e,
	}
}