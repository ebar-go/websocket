/**
 * @Author: Hongker
 * @Description:
 * @File:  interface
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

import (
	"github.com/ebar-go/websocket/epoll"
	cmap "github.com/orcaman/concurrent-map"
	"log"
	"net/http"
)

// Server websocket服务
type Server interface {
	// 处理请求
	HandleRequest(w http.ResponseWriter, r *http.Request)
	// 连接时触发
	HandleConnect(callback Callback)
	// 断开连接时触发
	HandleDisconnect(callback Callback)
	// 映射路由
	Route(uri string, handler Handler)
	// 关闭连接
	Close(conn Connection)
	// 广播
	Broadcast(response Response, ignores ...string)
	// 启动服务
	Start()
}



// NewServer 多进程server，相比epoll的单进程，降低了延迟
func NewServer(opts ...Option) Server {
	e, err := epoll.Create()
	if err != nil {
		log.Fatalf("create epoll:%v\n", err)
	}
	return &server{
		engine: newEngine(),
		connections: cmap.New(),
		epoller: e ,
		workers: newWorkerPool(opts...),
	}
}