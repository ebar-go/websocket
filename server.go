/**
 * @Author: Hongker
 * @Description:
 * @File:  server
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

import (
	"net/http"
	"sync"
)

type server struct {
	rmw sync.RWMutex
	engine *Engine
	connections map[string]Connection
	register    chan Connection
	unregister  chan Connection
}

func (s *server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := newConnection(w, r)
	if err != nil {
		// do something..
		return
	}

	s.registerConn(conn)
}

func (s *server) HandleConnect(f func(conn Connection)) {
	panic("implement me")
}

func (s *server) HandleDisconnect(f func(conn Connection)) {
	panic("implement me")
}

func (s *server) Route(uri string, handler func(ctx Context)) {
	s.engine.route(uri, handler)
}


func (s *server) Broadcast(response Response) {
	panic("implement me")
}

// Register register conn
func (srv *server) registerConn(conn Connection) {
	go func() {
		defer srv.Close(conn)

		srv.engine.listen(conn)
	}()
	srv.register <- conn
}

func (srv *server) Close(conn Connection)  {
	conn.close()
	srv.unregister <- conn
}

// Start
func (srv *server) Start() {
	if srv.engine.noRoute == nil {
		srv.engine.NoRoute(notFoundHandler)
	}
	go func() {
		for {
			select {
			case conn := <-srv.register:
				srv.connections[conn.ID()] = conn
			case conn := <-srv.unregister:
				if _, ok := srv.connections[conn.ID()]; ok {
					delete(srv.connections, conn.ID())
				}
			}
		}
	}()
}

func NewServer() Server {
	return &server{
		engine: &Engine{routers: map[string]Handler{}},
		connections: make(map[string]Connection),
		register:    make(chan Connection),
		unregister:  make(chan Connection),
	}
}