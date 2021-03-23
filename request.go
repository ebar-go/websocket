/**
 * @Author: Hongker
 * @Description:
 * @File:  request
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

type MessageRequest struct {
	RequestUri string `json:"uri"`
	RequestBody string `json:"body"`
}

func (r MessageRequest) Uri() string {
	return r.RequestUri
}

func (r MessageRequest) Body() []byte {
	return []byte(r.RequestBody)
}

