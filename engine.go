/**
 * @Author: Hongker
 * @Description:
 * @File:  engine
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

import (
	cmap "github.com/orcaman/concurrent-map"
	"net/http"
)

// Engine 路由引擎
type Engine struct {
	// 路由映射,未来考虑升级为支持restful模式
	routers cmap.ConcurrentMap
	// 404
	notFound Handler
}

// Handler handle connection request
type Handler func(ctx Context)

// route 设置路由映射
func (engine *Engine) route(uri string, handler Handler) {
	engine.routers.Set(uri, handler)
}

// Handle 执行路由
func (engine *Engine) handle(ctx Context) {
	// 获取路由映射的handler
	handler, exist := engine.routers.Get(ctx.RequestUri())
	if !exist {
		// 404
		engine.notFound(ctx)
		return
	}

	handler.(Handler)(ctx)
}

// NoRoute 设置404处理器
func (engine *Engine) NoRoute(handler Handler) {
	if handler == nil {
		return
	}
	engine.notFound = handler
}

// newEngine 实例
func newEngine() *Engine {
	return &Engine{
		notFound: func(ctx Context) {
			ctx.Error(http.StatusNotFound, "404 not found")
		},
		routers: cmap.New(),
	}
}
