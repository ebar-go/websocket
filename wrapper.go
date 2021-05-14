package websocket

import (
	"encoding/json"
	"log"
)

// ContextWrapper 增加其输出带statusCode的扩展功能
type ContextWrapper struct {
	ctx Context
}

func Wrap(ctx Context) ContextWrapper {
	return ContextWrapper{ctx: ctx}
}

// Success implement of Context
func (wrapper ContextWrapper) Success(data interface{}) {
	wrapper.write(0, "success", data)
}

// Error 错误信息
func (wrapper ContextWrapper) Error(code int, message string) {
	wrapper.write(code, message, nil)
}

// write 发送信息到客户端
func (wrapper ContextWrapper) write(code int, message string, data interface{}) {
	p := customResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}

	err := wrapper.ctx.WriteJson(p)
	if err != nil {
		log.Println("unable to write message:", err.Error())
	}
}


// customResponse 自定义响应内容
type customResponse struct {
	// 状态码,0为成功
	Code int `json:"code"`
	// 提示信息
	Message string `json:"message"`
	// 数据项
	Data interface{} `json:"data"`
}

// Byte implement of Response
func (response customResponse) Byte() []byte {
	b, _ := json.Marshal(response)
	return b
}