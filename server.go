/**
 * @Author: Hongker
 * @Description:
 * @File:  simpleServer
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"log"
)

// 基础server
type server struct {
	// socket连接
	connections cmap.ConcurrentMap
	// 路由引擎
	engine *Engine
	// 连接回调
	connectCallback func(conn Connection)
	// 注销回调
	disconnectCallback func(conn Connection)
}

func base() server {
	return server{engine: newEngine(), connections: cmap.New()}
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

func (srv *server) key(fd int) string {
	return fmt.Sprintf("idx:%d", fd)
}

func (srv *server) AddConnection(conn Connection) {
	srv.connections.Set(srv.key(conn.fd()), conn)
	// 注册回调
	if srv.connectCallback != nil {
		srv.connectCallback(conn)
	}
}

func (srv *server) RemoveConnection(conn Connection) {
	// 关闭socket
	defer conn.close()

	srv.connections.Remove(srv.key(conn.fd()))

	// 注销回调
	if srv.disconnectCallback != nil {
		srv.disconnectCallback(conn)
	}
}

func (srv *server) GetConnection(fd int) (Connection, bool) {
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