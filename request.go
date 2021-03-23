/**
 * @Author: Hongker
 * @Description:
 * @File:  request
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

type request struct {
	RequestUri string `json:"uri"`
	RequestBody string `json:"body"`
}

func (r request) Uri() string {
	return r.RequestUri
}

func (r request) Body() []byte {
	return []byte(r.RequestBody)
}

