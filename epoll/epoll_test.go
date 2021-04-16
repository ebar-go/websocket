package epoll

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T)  {
	s := make([]int, 4)
	fmt.Println(s, len(s), cap(s))
}
