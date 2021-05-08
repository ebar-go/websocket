/**
 * @Author: Hongker
 * @Description:
 * @File:  interface
 * @Version: 1.0.0
 * @Date: 2021/3/23 20:08
 */

package websocket

import (
	"github.com/ebar-go/websocket/epoll"
	cmap "github.com/orcaman/concurrent-map"
	"log"
)

// Context 上下文
type Context interface {
	// 获取header信息
	GetHeader(key string) string
	// 获取请求资源
	RequestUri() string
	// 通过json解析body
	BindJson(obj interface{}) error
	// 输出成功的数据
	Success(data interface{})
	// 错误信息
	Error(code int, message string)
}

// Data 数据项
type Data map[string]interface{}

// NewServer 多进程server，相比epoll的单进程，降低了延迟
func NewServer(opts ...Option) *Server {
	e, err := epoll.Create()
	if err != nil {
		log.Fatalf("failed to create epoll:%v\n", err)
	}

	// default option
	option := options{
		workers: 50,
		tasks:   100000,
	}
	for _, opt := range opts {
		opt.apply(&option)
	}

	return &Server{
		engine:      newEngine(),
		connections: cmap.New(),
		epoller:     e,
		workers:     newWorkerPool(option.workers, option.tasks),
	}
}
