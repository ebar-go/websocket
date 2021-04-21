/**
 * only websocket
 */
package main

import (
	"github.com/ebar-go/websocket"
	"log"
	"net/http"
)

func main() {
	ws := websocket.NewServer()

	// 监听连接创建事件
	ws.HandleConnect(func(conn websocket.Connection) {
		log.Printf("welcome: %s\n", conn.ID())
		conn.Write([]byte("hello"))
	})
	// 监听连接断开事件
	ws.HandleDisconnect(func(conn websocket.Connection) {
		log.Printf("goodbye: %s\n", conn.ID())
	})

	// 路由以及handler
	ws.Route("/index", func(ctx websocket.Context) {
		req := struct {
			Name string `json:"name"`
		}{}

		if err := ctx.BindJson(&req); err != nil {
			log.Println(err)
			ctx.Error(1001, "参数错误")
			return
		}
		ctx.Success(websocket.Data{
			"name": req.Name,
		})
	})

	// 启动
	ws.Start()

	// 绑定http服务
	http.ListenAndServe(":8085", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.HandleRequest(w, r)
	}))
}