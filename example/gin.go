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
	router := gin.Default()
	//ws := websocket.SampleServer() // 基于goroutine-per-conn实现的server
	//ws := websocket.EpollServer() // 基于epoll实现的server
	ws := websocket.WorkerPoolServer(50, 100000) // 基于workerPool实现的epollServer
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
