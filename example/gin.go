/**
 * @Author: Hongker
 * @Description:
 * @File:  gin
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:55
 */

package main

import (
	"github.com/ebar-go/websocket"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 实例化web服务
	router := gin.Default()
	ws := websocket.NewServer(
		websocket.WithWorkerNumber(1),  // 设置worker数量，可选，默认为50
		websocket.WithTaskNumber(2), // 设置task数量，可选，默认为100000
		) // 基于workerPool实现的epollServer
	// 用于创建websocket连接
	router.GET("/ws", func(ctx *gin.Context) {
		ws.HandleRequest(ctx.Writer, ctx.Request)
	})

	// 广播
	router.GET("/broadcast", func(ctx *gin.Context) {
		ws.Broadcast(websocket.Response{
			Code:    0,
			Message: "test",
			Data:    "hello,websocket",
		})
		ctx.JSON(200, gin.H{"hello": "world"})
	})
	// 监听连接创建事件
	ws.HandleConnect(func(conn websocket.Connection) {
		log.Printf("welcome: %s\n", conn.ID())
	})
	// 监听连接断开事件
	ws.HandleDisconnect(func(conn websocket.Connection) {
		log.Printf("goodbye: %s\n", conn.ID())
	})

	// 路由以及handler
	ws.Route("/index", func(ctx websocket.Context) {
		ctx.Success("hello,world")
	})

	// 启动
	ws.Start()

	router.Run(":8091")
}
