/**
 * @Author: Hongker
 * @Description:
 * @File:  request
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

// Request ws请求
type Request interface {
	// 获取请求路由
	Uri() string
	// 获取请求内容
	Body() []byte
}
// request 自定义Request
type request struct {
	RequestUri string `json:"uri"`
	RequestBody string `json:"body"`
}

// Uri implement of Request
func (r request) Uri() string {
	return r.RequestUri
}

// Body implement of Request
func (r request) Body() []byte {
	return []byte(r.RequestBody)
}

