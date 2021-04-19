/**
 * @Author: Hongker
 * @Description:
 * @File:  request
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

// request 自定义Request
type Request struct {
	Uri string `json:"uri"`
	Body string `json:"body"`
}

func (req Request) body() []byte {
	return []byte(req.Body)
}

