# websocket
像web服务一样去处理websocket

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
	ws := websocket.Default()
	router.GET("/ws", func(ctx *gin.Context) {
		ws.HandleRequest(ctx.Writer, ctx.Request)
	})
	ws.HandleConnect(func(conn websocket.Connection) {
		log.Printf("welcome: %s\n", conn.ID())
	})
	ws.HandleDisconnect(func(conn websocket.Connection) {
		log.Printf("goodbye: %s\n", conn.ID())
	})

	ws.Route("/index", func(ctx websocket.Context) {
		ctx.Success("hello,world")
	})

	ws.Start()

	router.Run(":8081")
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