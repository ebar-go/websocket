package websocket

import (
	"github.com/ebar-go/websocket/tree"
	"path"
)

type Router struct {
	// 前缀
	prefix string
	// 路由树
	tree *tree.RadixTree
}

func (router *Router) print() {
	router.tree.Print("")
}

// withPrefix 拼接前缀，获取到的是完整url地址
func (router *Router) withPrefix(uri string) string{
	return path.Join(router.prefix, uri)
}

// Group 分组
func (router *Router) Group(path string) *Router {
	path = router.withPrefix(path)
	child := &Router{prefix: path, tree: router.tree}
	return child
}

// Use 中间件
func (router *Router) Use() {
	// TODO
}
// Route 路由匹配
func (router *Router) Route(path string, handler Handler)  {
	router.tree.Insert(router.withPrefix(path), handler)
}

// Get 获取handler
func (router *Router) Get(path string) (Handler, bool){
	val := router.tree.GetValue(path)
	if val == nil {
		return nil, false
	}
	handler, ok := router.tree.GetValue(path).(Handler)
	return handler, ok
}

func NewRouter() *Router {
	return &Router{
		prefix: "/",
		tree:   tree.NewRadixTree(),
	}
}