/**
 * @Author: Hongker
 * @Description:
 * @File:  simpleServer
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

// 基础server
type server struct {
	// 路由引擎
	engine *Engine
	// 连接回调
	connectCallback func(conn Connection)
	// 注销回调
	disconnectCallback func(conn Connection)
}

func base() server {
	return server{engine: newEngine()}
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

