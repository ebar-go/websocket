package websocket

import (
	"github.com/ebar-go/websocket/tree"
	"path"
)

// Router 路由接口，面向接口开发
type Router interface {
	// Route 设置路由
	Route(path string, handler Handler)
	// Get 获取handler
	Get(path string) (Handler, bool)
	// Group 生成分组路由
	Group(path string) Router
}

// radixRouter 路由实现
type radixRouter struct {
	// 前缀
	prefix string
	// 路由树
	tree *tree.RadixTree
}

func (router *radixRouter) print() {
	router.tree.Print("")
}

// withPrefix 拼接前缀，获取到的是完整url地址
func (router *radixRouter) withPrefix(uri string) string{
	return path.Join(router.prefix, uri)
}

// Group 新生成一个前缀为path的路由
func (router *radixRouter) Group(path string) Router {
	path = router.withPrefix(path)
	child := &radixRouter{prefix: path, tree: router.tree}
	return child
}

// Use 中间件
func (router *radixRouter) Use() {
	// TODO
}
// Route 路由匹配
func (router *radixRouter) Route(path string, handler Handler)  {
	router.tree.Insert(router.withPrefix(path), handler)
}

// Get 获取handler
func (router *radixRouter) Get(path string) (Handler, bool){
	val := router.tree.GetValue(path)
	if val == nil {
		return nil, false
	}
	handler, ok := router.tree.GetValue(path).(Handler)
	return handler, ok
}

// NewRouter 获取路由实例
func NewRouter() Router {
	return &radixRouter{
		prefix: "/",
		tree:   tree.NewRadixTree(),
	}
}