package router

import (
	"strings"
	"sync"
)

// RadixTree 空间优化前缀树
type RadixTree struct {
	// 锁
	mu sync.RWMutex
	// 子节点
	children map[string]*RadixTree
	// 是否为字符
	isWord   bool
	// 值
	val interface{}
}

// comparePrefix 遍历查询共同前缀，返回下标
func (rt *RadixTree) comparePrefix(s1, s2 string) int {
	// 获取短的字符的长度
	max := len(s1)
	if l := len(s2); l < max {
		max = l
	}

	var i int
	for i = 0; i < max; i++ {
		if s1[i] != s2[i] {
			break
		}
	}

	return i
}

// Insert 插入
func (rt *RadixTree) Insert(word string, val interface{}) {
	// lock for concurrent map
	rt.mu.Lock()
	defer rt.mu.Unlock()
	// 遍历子节点
	for s, child := range rt.children {
		// 1.判断目标字符串与子节点字符串是否相等，则直接标记isWord即可
		if s == word {
			child.isWord = true
			child.val = val
			continue
		}

		// 对比前缀
		idx := rt.comparePrefix(s, word)
		if idx == 0 { // 如果没有找到共同前缀，则继续找
			continue
		}

		// 发现共同前缀
		samePrefix := word[0:idx] // s[0:idx]一样

		// 如果当前子节点就是共同前缀,则直接插入剩余部分
		if s == samePrefix {
			child.Insert(word[idx:], val)
			return
		}

		// 当有共同前缀，则需要删除原来的节点，创建新的共同前置节点
		delete(rt.children, s)
		node := newRadixTreeNode(false)
		// 插入原节点的剩余部分
		node.Insert(s[idx:], child.val)
		// 如果插入字符串和前缀一致，则标记新节点为字符串节点
		if word == samePrefix {
			node.isWord = true
		}else { // 只是部分相同，则插入剩余部分
			node.Insert(word[idx:], val)
		}
		// 绑定新的子节点
		rt.children[samePrefix] = node
		return

	}

	// 遍历完子节点，依然没有找到拥有共同前缀的时候，就新插入一个节点,比如 hello和world没有共同前缀，则插入新的world节点
	node := newRadixTreeNode(true)
	node.val = val
	rt.children[word] = node
}

// Search 查询
func (rt *RadixTree) Search(word string) bool {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	for s, child := range rt.children {
		if s == word { // 检查子节点是否存在
			return child.isWord // 且同时判断是否为字符串
		}

		if strings.HasPrefix(word, s) { // 如果当前子节点为目标字符串的前缀，则递归查询剩余字符串
			return child.Search(word[len(s):])
		}
	}

	return false
}

// Delete 删除
func (rt *RadixTree) Delete(word string) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	for s, child := range rt.children {
		if s == word { // 检查子节点是否存在
			if len(child.children) == 0 {
				// 如果子节点没有更后面的节点了，直接删除该子节点
				delete(rt.children, s)
				return
			}
			child.isWord = false // 存在直接设为false
			child.val = nil
		}

		// 如果只是部分前缀，则继续遍历删除
		if strings.HasPrefix(word, s) {
			child.Delete(word[len(s):])
		}
	}
}

func (rt *RadixTree) GetValue(word string) interface{} {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	for s, child := range rt.children {
		if s == word { // 检查子节点是否存在
			return child.val
		}

		if strings.HasPrefix(word, s) { // 如果当前子节点为目标字符串的前缀，则递归查询剩余字符串
			return child.GetValue(word[len(s):])
		}
	}

	return nil
}

// newRadixTreeNode 生成新节点
func newRadixTreeNode(isWold bool) *RadixTree {
	node := new(RadixTree)
	node.children = make(map[string]*RadixTree, 26)
	node.isWord = isWold
	return node
}

// NewRadixTree
func NewRadixTree() *RadixTree {
	return newRadixTreeNode(false)
}