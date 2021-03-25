/**
 * @Author: Hongker
 * @Description:
 * @File:  engine
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

import (
	"net/http"
	"sync"
)
// Engine 路由引擎
type Engine struct {
	// 读写锁
	rmw sync.RWMutex
	// 路由映射,未来考虑升级为支持restful模式
	routers map[string]Handler
	// 404
	noRoute Handler
}
// route 设置路由映射
func (engine *Engine) route(uri string, handler Handler) {
	// 加锁，避免并发
	engine.rmw.Lock()
	defer engine.rmw.Unlock()
	engine.routers[uri] = handler
}
// Handle 执行路由
func (engine *Engine) handle(ctx Context) {
	// 加锁，避免并发
	engine.rmw.RLock()
	defer engine.rmw.RUnlock()

	// 获取路由映射的handler
	handler, ok := engine.routers[ctx.Request().Uri()]
	if !ok {
		// 404
		engine.noRoute(ctx)
		return
	}

	handler(ctx)
}
// NoRoute 设置404处理器
func (engine *Engine) NoRoute(handler Handler) {
	if handler == nil {
		return
	}
	engine.noRoute = handler
}

// listen 监听连接
func (engine *Engine) listen(conn Connection) {
	for {
		ctx, err := conn.context()
		if err != nil {
			break
		}

		engine.handle(ctx)
	}
}
// notFoundHandler 默认的handler
func notFoundHandler(ctx Context)  {
	ctx.Render(&response{
		Code:    http.StatusNotFound,
		Message: "404 not found",
		Data:    nil,
	})
}

func engine() *Engine {
	return &Engine{noRoute: notFoundHandler, routers: map[string]Handler{}}
}