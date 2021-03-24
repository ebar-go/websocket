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
	ws := websocket.EpollServer()
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

	router.Run(":8091")
}
