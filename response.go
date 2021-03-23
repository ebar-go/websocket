/**
 * @Author: Hongker
 * @Description:
 * @File:  response
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

import "encoding/json"

// Response websocket的响应数据结构
type Response interface {
	// 输出为[]byte
	Byte() []byte
}

// response 自定义响应内容
type response struct {
	// 状态码,0为成功
	Code int `json:"code"`
	// 提示信息
	Message string `json:"message"`
	// 数据项
	Data interface{} `json:"data"`
}

// Byte implement of Response
func (r response) Byte() []byte {
	b, _ := json.Marshal(r)
	return b
}

