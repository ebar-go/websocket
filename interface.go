/**
 * @Author: Hongker
 * @Description:
 * @File:  interface
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

import "net/http"

// Server 接口
type Server interface {
	// HandleRequest 处理请求
	HandleRequest(w http.ResponseWriter, r *http.Request)
	// Route 设置路由映射
	Route(uri string, handler Handler)
	// Close 主动关闭连接
	Close(conn Connection)

	// TODO 实现通过连接的唯一ID来关闭连接

	// Broadcast 广播
	Broadcast(message []byte, ignores ...string)
	// Start 启动服务
	Start()
	// Group 分组路由
	Group(uri string) Router
}

// Context 上下文
type Context interface {
	// GetHeader 获取header信息
	GetHeader(key string) string
	// RequestUri 获取请求资源
	RequestUri() string
	// BindJson 通过json解析body
	BindJson(obj interface{}) error
	// WriteJson 输出json
	WriteJson(obj interface{}) error
	// WriteMessage 输出byte
	WriteMessage(message []byte) error
	// WriteString 输出字符串
	WriteString(message string) error
}

// Data 数据项
type Data map[string]interface{}

