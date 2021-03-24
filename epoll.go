/**
 * @Author: Hongker
 * @Description:
 * @File:  epoll
 * @Version: 1.0.0
 * @Date: 2021/3/24 22:33
 */

package websocket

import (
	"github.com/gorilla/websocket"
	"golang.org/x/sys/unix"
	"reflect"
	"sync"
	"syscall"
)

type epoll struct {
	fd int
	connections map[int]*websocket.Conn
	lock sync.RWMutex
}

func MkEpoll() (*epoll, error) {
	fd, err := unix.EpollCreate(0)
	if err != nil {
		return nil, err
	}

	return &epoll{
		fd:          fd,
		connections: make(map[int]*websocket.Conn),
		lock:        sync.RWMutex{},
	}, nil
}

func websocketFD(conn *websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn).FieldByName("conn").Elem())
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func (e *epoll) Add(conn *websocket.Conn) error {
	fd := websocketFD(conn)
	err := unix.EpollCtl(e.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()
	e.connections[fd] = conn
	return nil
}

func (e *epoll) Remove(conn *websocket.Conn) error{
	fd := websocketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.connections, fd)

	return nil
}

func (e *epoll) Wait() ([]*websocket.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return nil, err
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	var connections []*websocket.Conn
	for i := 0; i < n; i++ {
		conn := e.connections[int(events[i].Fd)]
		connections = append(connections, conn)
	}
	return connections, nil
}