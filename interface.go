/**
 * @Author: Hongker
 * @Description:
 * @File:  interface
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

import "net/http"

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

type Handler func(ctx Context)