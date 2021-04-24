package router

import "testing"

func TestRadixTree_Print(t *testing.T) {
	tree := NewRadixTree()
	tree.Insert("sleep", 1)
	tree.Insert("son", 2)
	tree.Insert("sex", 3)
	tree.Print("")
}
