package context

import "encoding/json"


// Response 自定义响应内容
type Response struct {
	// 状态码,0为成功
	Code int `json:"code"`
	// 提示信息
	Message string `json:"message"`
	// 数据项
	Data interface{} `json:"data"`
}

// Byte implement of Response
func (r Response) Byte() []byte {
	b, _ := json.Marshal(r)
	return b
}
