/**
 * @Author: Hongker
 * @Description:
 * @File:  context
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

import (
	"encoding/json"
	"log"
)

// Context 上下文
type Context interface {
	// 获取请求资源
	Uri() string
	// 通过json解析body
	BindJson(obj interface{}) error
	// 输出response
	Render(response Response)
	// 输出成功的数据
	Success(data interface{})
}

// context 自定义context
type context struct {
	// 请求体
	request Request
	// 当前连接
	conn Connection
}
// Uri return request uri
func (ctx *context) Uri() string {
	return ctx.request.Uri()
}
// BindJson implement of Context
func (ctx *context) BindJson(obj interface{}) error {
	return json.Unmarshal(ctx.request.Body(), obj)
}
// Render implement of Context
func (ctx *context) Render(response Response) {
	if err := ctx.conn.write(response.Byte()); err != nil {
		log.Println("unable to write message:", err.Error())
	}
}
// Success implement of Context
func (ctx *context) Success(data interface{}) {
	ctx.Render(&response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func (ctx *context) Error(code int, message string) {
	ctx.Render(&response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}



