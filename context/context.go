package context

import (
	"github.com/gorilla/websocket"
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

// RequestUri 获取请求的uri
func (ctx *Context) RequestUri() string {
	return ctx.request.Uri
}

// BindJson implement of Context
func (ctx *Context) BindJson(obj interface{}) error {
	return ctx.request.Unmarshal(obj)
}

func (ctx *Context) WriteJson(obj interface{}) error {
	return ctx.conn.WriteJSON(obj)
}

func (ctx *Context) WriteMessage(message []byte) error {
	return ctx.conn.WriteMessage(websocket.TextMessage, message)
}


func (ctx *Context) WriteString(message string) error {
	return ctx.conn.WriteMessage(websocket.TextMessage, []byte(message))
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
