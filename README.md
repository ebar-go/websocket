# websocket
- 基于Epoll和WorkerPool实现的高性能websocket框架。 
- 提供路由树模式，让开发者像开发http接口一样方便的去开发websocket应用。
- 支持创建连接与断开连接的回调事件。

## 安装
```
go get github.com/ebar-go/websocket
```

## 初始化
```go
package main


import (
	"github.com/ebar-go/websocket"
	"log"
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
}

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
	ws := websocket.NewBuilder().Build()
	
	ws.Route("home", func(ctx websocket.Context) {
		ctx.WriteString("welcome")
	})
	// 支持路由分组
	userGroup := ws.Group("user")
	{
		// 请求uri为: /user/list
		userGroup.Route("list", func(ctx websocket.Context) {
			// 直接输出字符串
			ctx.WriteString("this is user list api")
		})
		// 请求的uri为/user/create
		userGroup.Route("create", func(ctx websocket.Context) {
			// 输出json
			ctx.WriteJson(map[string]interface{}{"code":0})
		})
	}
	
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
	ws := websocket.NewBuilder().Build()
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
> {"uri":"/home"}
< welcome
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
func main() {
    // 路由以及handler,请求uri为: /index
	ws.Route("index", func(ctx websocket.Context) {
		req := struct {
			Name string `json:"name"`
		}{}
		// 解析参数
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
}
```

## 技术细节
- 使用epoll模型管理连接，减少了goroutine的创建
- 使用worker pool提高吞吐率
- 使用了建造者模式、装饰器模式

## TODO
- 压力测试
- 中间件