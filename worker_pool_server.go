package websocket

import (
	"log"
	"net/http"
)

type workerPoolServer struct {
	server

	epoller *epoll

	workerPool *pool
}

func (srv *workerPoolServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// 获取socket连接
	conn, err := newConnection(w, r)
	if err != nil {
		// do something..
		return
	}
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



func (srv *workerPoolServer) Close(conn Connection) {
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

func (srv *workerPoolServer) Broadcast(response Response, ignores ...string) {
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

func (srv *workerPoolServer) Start() {
	// 设置默认的404路由
	if srv.engine.noRoute == nil {
		srv.engine.NoRoute(notFoundHandler)
	}
	srv.workerPool.start(srv.engine)

	go func() {
		defer srv.workerPool.Close()
		for {
			connections, err := srv.epoller.Wait()
			if err != nil {
				log.Printf("Failed to epoll wait %v", err)
				continue
			}
			for _, conn := range connections {
				if conn == nil {
					break
				}
				srv.workerPool.addTask(conn)
			}
		}

	}()
}

// WorkerPoolServer 多进程server，相比epoll的单进程，降低了延迟
//
// workers 进程数
//
// maxConnections 最大并发连接数
func WorkerPoolServer(workers , maxConnections int) Server {
	epoller, err := MkEpoll()
	if err != nil {
		log.Fatalf("create epoll:%v\n", err)
	}
	return &workerPoolServer{
		server: base(),
		epoller: epoller ,
		workerPool: newPool(workers, maxConnections),
	}
}
