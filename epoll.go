/**
 * @Author: Hongker
 * @Description:
 * @File:  epoll
 * @Version: 1.0.0
 * @Date: 2021/3/24 22:33
 */

package websocket

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"golang.org/x/sys/unix"
	"syscall"
)
// epoll struct
type epoll struct {
	// 句柄
	fd int
	// socket 连接
	connections cmap.ConcurrentMap
}

// MkEpoll return epoll instance
func MkEpoll() (*epoll, error) {
	fd, err := unix.EpollCreate(1)
	if err != nil {
		return nil, err
	}

	return &epoll{
		fd:          fd,
		connections: cmap.New(),
	}, nil
}

// key format map index
func (e *epoll) key(fd int) string {
	return fmt.Sprintf("f%d", fd)
}


func (e *epoll) Add(conn Connection) error {
	fd := conn.fd()
	err := unix.EpollCtl(e.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}

	e.connections.Set(e.key(fd), conn)
	return nil
}

func (e *epoll) Remove(conn Connection) error{
	fd := conn.fd()
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	e.connections.Remove(e.key(fd))

	return nil
}

func (e *epoll) Wait() ([]Connection, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return nil, err
	}
	connections := make([]Connection, 0, n)
	for i := 0; i < n; i++ {
		conn, exist := e.connections.Get(e.key(int(events[i].Fd)))
		if exist {
			connections[i] = conn.(Connection)
		}

	}
	return connections, nil
}