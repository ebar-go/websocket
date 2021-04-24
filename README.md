# websocket
- 基于Epoll和WorkerPool实现的高性能websocket框架。 
- 提供路由树模式，让开发者像开发http接口一样方便的去开发websocket应用。
- 支持创建连接与断开连接的回调事件。

## 安装
```
go get github.com/ebar-go/websocket
```

## Demo
- only websocket
```go
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
	// 支持路由分组
	userGroup := ws.Group("user")
	{
		// 请求uri为: /user/list
		userGroup.Route("list", func(ctx websocket.Context) {
			ctx.Success("this is user list api")
		})
		// 请求的uri为/user/create
		userGroup.Route("create", func(ctx websocket.Context) {
			ctx.Success("this is user create api")
		})
	}
	// 路由以及handler,请求uri为: /index
	ws.Route("index", func(ctx websocket.Context) {
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
```

- gin + websocket
```go
package main

import (
	"github.com/ebar-go/websocket"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	ws := websocket.NewServer()
	// 用于创建websocket连接
	router.GET("/ws", func(ctx *gin.Context) {
		ws.HandleRequest(ctx.Writer, ctx.Request)
	})
	// other demo
	ws.Start()

	router.Run(":8085")
}
```

通过`wscat`去连接websocket:
```
wscat -c ws://127.0.0.1:8081/ws
> {"uri":"/index"}
< {"code":0,"message":"success","data":"hello,world"}
> {"uri":"/home"}
< {"code":404,"message":"404 not found","data":null}
```

## 请求参数
- uri: 路由
- body: 数据内容
```json
{"uri":"/index","body": {"name": "websocket"}}
```
- 获取参数
```
req := struct {
    Name string `json:"name"`
}{}
// 通过BindJson解析数据
if err := ctx.BindJson(&req); err != nil {
    ctx.Error(1001, "参数错误")
    return
}
ctx.Success(websocket.Data{
    "name": req.Name,
})
```
## 压力测试
TODO