package context

import "github.com/gorilla/websocket"

// ContextBuilder context的建造器
type ContextBuilder struct {
	// conn 连接
	conn *websocket.Conn

	// 其他参数
}

func NewContextBuilder(conn *websocket.Conn) *ContextBuilder  {
	return &ContextBuilder{
		conn : conn,
	}
}

func (builder *ContextBuilder) Build() *Context {
	return &Context{
		conn:     builder.conn,
		request:  Request{},
		index:    0,
		handlers: nil,
	}
}
