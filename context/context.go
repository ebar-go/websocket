package context

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Context struct {
	conn *websocket.Conn
	request Request
}


// RequestUri
func (ctx *Context) RequestUri() string {
	return ctx.request.Uri
}

// BindJson implement of Context
func (ctx *Context) BindJson(obj interface{}) error {
	// marshal json
	b, err := ctx.request.jsonMarshal()
	if err != nil {
		return err
	}
	// unmarshal json
	return json.Unmarshal(b, obj)
}

// Success implement of Context
func (ctx *Context) Success(data interface{}) {
	ctx.write(0, "success", data)
}


func (ctx *Context) Error(code int, message string) {
	ctx.write(code, message, nil)
}

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

func NewContext(conn *websocket.Conn) *Context {
	return &Context{
		conn: conn,
	}
}

func (ctx *Context) Read() error {
	_, message, err := ctx.conn.ReadMessage()
	if err != nil {
		return  err
	}

	req, err := NewRequest(message)
	if err != nil {
		return err
	}
	ctx.request = req
	return nil
}