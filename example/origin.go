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
	ws := websocket.NewBuilder().
		SetWorkerNumber(50). // 设置并发worker数量，影响吞吐率，默认为50，可根据实际测试数据来评估，可选
		SetMaxConnectionNumber(100000). // 设置最大连接数，默认为100000，可根据最大并发量来评估，可选
		SetConnectCallback(func(conn websocket.Connection) { // 设置连接创建回调，可选
			log.Printf("welcome: %s\n", conn.ID())
		}).SetDisconnectCallback(func(conn websocket.Connection) { // 设置连接断开回调，可选
			log.Printf("goodbye: %s\n", conn.ID())
		}).Build()

	// 路由分组
	userGroup := ws.Group("user")
	{
		// 请求uri为: /user/list
		userGroup.Route("list", func(ctx websocket.Context) {
			_ = ctx.WriteString("this is user list api")
		})
		// 请求的uri为/user/create
		userGroup.Route("create", func(ctx websocket.Context) {
			_ = ctx.WriteString("this is user create api")
		})
	}

	// 请求uri为: /index
	ws.Route("index", func(ctx websocket.Context) {
		req := struct {
			Name string `json:"name"`
		}{}

		if err := ctx.BindJson(&req); err != nil {
			log.Println(err)
			websocket.Wrap(ctx).Error(1001, "参数错误")
			return
		}

		// 使用ContextWrapper输出带code的响应内容
		websocket.Wrap(ctx).Success(websocket.Data{
			"name": req.Name,
		})
	})

	// 启动
	ws.Start()

	// 绑定http服务
	_ = http.ListenAndServe(":8085", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.HandleRequest(w, r)
	}))
}
