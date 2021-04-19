# websocket
- 基于Epoll和WorkerPool实现的高性能websocket框架。 
- 提供路由模式，让开发者像开发http接口一样方便的去开发websocket应用。


## 安装
```
go get github.com/ebar-go/ego
```
## Demo
```go

package main

import (
	"github.com/ebar-go/websocket"
	"github.com/gin-gonic/gin"
	"log"
)


func main() {
	router := gin.Default()
	//ws := websocket.SampleServer() // 基于goroutine-per-conn实现的server
	//ws := websocket.EpollServer() // 基于epoll实现的server
	ws := websocket.NewServer(
		websocket.WithWorkerNumber(1),  // 设置worker数量，可选，默认为50
		websocket.WithTaskNumber(2), // 设置task数量，可选，默认为100000
	) // 基于workerPool实现的epollServer
	// 用于创建websocket连接
	router.GET("/ws", func(ctx *gin.Context) {
		ws.HandleRequest(ctx.Writer, ctx.Request)
	})
	// 监听连接创建事件
	ws.HandleConnect(func(conn websocket.Connection) {
		log.Printf("welcome: %s\n", conn.ID())
	})
	// 监听连接断开事件
	ws.HandleDisconnect(func(conn websocket.Connection) {
		log.Printf("goodbye: %s\n", conn.ID())
	})

	// 路由及处理函数
	ws.Route("/index", func(ctx websocket.Context) {
		ctx.Success("hello,world")
	})

	ws.Start()

	router.Run(":8091")
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
```go
package main

import (
	"github.com/ebar-go/websocket"
)

func main() {
    ws := websocket.NewServer()
    // ...
	// 路由以及handler
	ws.Route("/index", func(ctx websocket.Context) {
		// 定义结构体
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
	})
}
```
## 压力测试
TODO