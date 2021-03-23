/**
 * @Author: Hongker
 * @Description:
 * @File:  context
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:07
 */

package websocket

import "encoding/json"

// Context 上下文
type Context interface {
	// 获取请求属性
	Request() Request
	// 通过json解析body
	BindJson(obj interface{}) error
	// 输出response
	Render(response Response)
	// 输出成功的数据
	Success(data interface{})
}

// context 自定义context
type context struct {
	request Request
	conn Connection
}

func (ctx *context) Request() Request {
	return ctx.request
}

func (ctx *context) BindJson(obj interface{}) error {
	return json.Unmarshal(ctx.request.Body(), obj)
}

func (ctx *context) Render(response Response) {
	_ = ctx.conn.write(response.Byte())
}

func (ctx *context) Success(data interface{}) {
	ctx.Render(&response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}



