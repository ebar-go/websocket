/**
 * 基于epoll和worker pool模式，实现可支持高连接数与高性能的websocket框架
 * epoll: 管理websocket连接的文件标识符，主要包含：添加，删除，获取活跃分子。
 * 	避免因一个连接一个goroutine的浪费内存问题(1个g需要2~8kb的内存，如果是百万连接，则会消耗几十个G的内存)
 *
 * worker pool: 通过多worker模式并发处理websocket请求，提高吞吐率。
 */

package websocket

import (
	"fmt"
	"github.com/ebar-go/websocket/epoll"
	cmap "github.com/orcaman/concurrent-map"
	"log"
	"net/http"
)


// server
type server struct {
	// socket连接,通过concurrent map,保证并发安全，同时提高性能
	connections cmap.ConcurrentMap

	// 路由引擎.处理请求
	engine *Engine

	// 连接回调
	connectCallback Callback

	// 注销回调
	disconnectCallback Callback

	// epoll
	epoller epoll.Epoll

	// worker pool
	workers *workerPool
}

// HandleRequest 处理websocket请求，主要是注册socket连接
func (srv *server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// 获取socket连接
	conn, err := newConnection(w, r)
	if err != nil {
		// do something..
		return
	}

	// 通过epoll模型来管理websocket连接
	if err := srv.epoller.Add(conn.fd()); err != nil {
		log.Println("unable to add connection:", err.Error())
		_ = conn.close()
		return
	}

	// 注册连接
	srv.registerConn(conn)
}

// Close connection
func (srv *server) Close(conn Connection) {
	// remove socket in epoll model
	if err := srv.epoller.Remove(conn.fd()); err != nil {
		log.Println("unable to remove conn:", err.Error())
	}
	// unregister connection
	srv.unregisterConn(conn)
}

// HandleConnect 建立连接时的回调
func (srv *server) HandleConnect(callback Callback) {
	srv.connectCallback = callback
}

// HandleDisconnect 断开连接时的回调
func (srv *server) HandleDisconnect(callback Callback) {
	srv.disconnectCallback = callback
}
// Route 绑定路由
func (srv *server) Route(uri string, handler Handler) {
	srv.engine.route(uri, handler)
}

// unique key
func (srv *server) key(fd int) string {
	return fmt.Sprintf("idx:%d", fd)
}
// registerConn add connection to map
func (srv *server) registerConn(conn Connection) {
	srv.connections.Set(srv.key(conn.fd()), conn)
	// 注册回调
	if srv.connectCallback != nil {
		srv.connectCallback(conn)
	}
}

// unregisterConn remove and close connection
func (srv *server) unregisterConn(conn Connection) {
	// 关闭socket
	defer conn.close()

	srv.connections.Remove(srv.key(conn.fd()))

	// 注销回调
	if srv.disconnectCallback != nil {
		srv.disconnectCallback(conn)
	}
}

// 通过文件标识符获取连接
func (srv *server) getConnection(fd int) (Connection, bool) {
	v, exist := srv.connections.Get(srv.key(fd))
	if !exist {
		return nil, false
	}
	return v.(Connection), true
}


// Broadcast implement of Server
func (srv *server) Broadcast(response Response, ignores ...string) {
	if len(ignores) == 0 {
		// not ignore
		srv.connections.IterCb(func(key string, v interface{}) {
			conn := v.(Connection)
			if err := conn.write(response.Byte()); err != nil {
				log.Printf("write to [%s]: %v", conn.ID(), err)
			}
		})

		return
	}

	// has ignores
	srv.connections.IterCb(func(key string, v interface{}) {
		conn := v.(Connection)
		for _, ignore := range ignores {
			if ignore == conn.ID() {
				return
			}
		}
		if err := conn.write(response.Byte()); err != nil {
			log.Printf("write to [%s]: %v", conn.ID(), err)
		}
	})
}

// Start
func (srv *server) Start() {
	log.Println("websocket serving..")
	// 给workers指定job
	srv.workers.job = func(fd int) {
		// 通过文件标识符，获取到socket连接
		conn , exist := srv.getConnection(fd)
		if !exist {
			return
		}

		ctx, err := conn.context()
		if err != nil {
			//srv.Close(conn)
			return
		}
		srv.engine.handle(ctx)
	}

	// 开始工作
	srv.workers.start()


	go func() {
		// 线程结束时，停止工作
		defer srv.workers.stop()
		for {
			// 通过wait方法获取到epoll管理的活跃socket连接
			fds, err := srv.epoller.Wait()
			if err != nil {
				log.Println("unable to get active socket connection from epoll:", err)
				continue
			}

			// 将连接分配给worker
			for _, fd := range fds {
				srv.workers.addTask(fd)
			}
		}

	}()

}
