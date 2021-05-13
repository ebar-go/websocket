package context

import (
	"github.com/gorilla/websocket"
	"log"
	"math"
)

const abortIndex = math.MaxInt8 / 2

// Context 上下文
type Context struct {
	// conn 连接
	conn *websocket.Conn
	// request 请求
	request Request
	// handler index
	index    uint8
	handlers []func()
}

// GetHeader 获取header信息
func (ctx *Context) GetHeader(key string) string {
	if ctx.request.Header == nil {
		return ""
	}

	return ctx.request.Header[key]
}

// RequestUri
func (ctx *Context) RequestUri() string {
	return ctx.request.Uri
}

// BindJson implement of Context
func (ctx *Context) BindJson(obj interface{}) error {
	return ctx.request.Unmarshal(obj)
}

// Success implement of Context
func (ctx *Context) Success(data interface{}) {
	ctx.write(0, "success", data)
}

// Error 错误信息
func (ctx *Context) Error(code int, message string) {
	ctx.write(code, message, nil)
}

// write 发送信息到客户端
func (ctx *Context) write(code int, message string, data interface{}) {
	p := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}.Byte()

	err := ctx.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		log.Println("unable to write message:", err.Error())
	}
}

// Read 读取socket信息
func (ctx *Context) Read() error {
	_, message, err := ctx.conn.ReadMessage()
	if err != nil {
		return err
	}

	req, err := newRequest(message)
	if err != nil {
		return err
	}
	ctx.request = req
	return nil
}

// next 继续执行,暂未实现
func (ctx *Context) next() {
	ctx.index++
	for ctx.index < uint8(len(ctx.handlers)) {
		ctx.handlers[ctx.index]()
		ctx.index++
	}
}

func (ctx *Context) isAborted() bool {
	return ctx.index >= abortIndex
}

// abort 中断，暂未实现
func (ctx *Context) abort() {
	ctx.index = abortIndex
}
