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

type Engine struct {
	rmw sync.RWMutex
	routers map[string]Handler
	noRoute Handler
}

func (engine *Engine) route(uri string, handler Handler) {
	engine.rmw.Lock()
	defer engine.rmw.Unlock()
	engine.routers[uri] = handler
}

func (engine *Engine) Run(ctx Context) {
	handler, ok := engine.routers[ctx.Request().Uri()]
	if !ok {
		// 404
		engine.noRoute(ctx)
		return
	}

	handler(ctx)
}

func (engine *Engine) NoRoute(handler Handler) {
	engine.noRoute = handler
}

// listen 监听连接
func (engine *Engine) listen(conn Connection) {
	for {
		ctx, err := conn.context()
		if err != nil {
			break
		}

		engine.Run(ctx)
	}
}

func notFoundHandler(ctx Context)  {
	ctx.Render(&response{
		Code:    http.StatusNotFound,
		Message: "404 not found",
		Data:    nil,
	})
}