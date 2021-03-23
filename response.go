/**
 * @Author: Hongker
 * @Description:
 * @File:  response
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

import "encoding/json"

type response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func (r response) Byte() []byte {
	b, _ := json.Marshal(r)
	return b
}

