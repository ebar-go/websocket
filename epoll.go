/**
 * @Author: Hongker
 * @Description:
 * @File:  epoll
 * @Version: 1.0.0
 * @Date: 2021/3/24 22:33
 */

package websocket

import (
	"net"
	"sync"
)

type epoll struct {
	fd int
	connections map[int]net.Conn
	lock sync.RWMutex
}

