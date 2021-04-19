package epoll

import (
	"golang.org/x/sys/unix"
	"syscall"
)

// Epoll model
type Epoll interface {
	// add fd
	Add(fd int) error
	// remove fd
	Remove(fd int) error
	// wait active fd
	Wait() ([]int, error)
}

// epollImpl is implement of Epoll
type epollImpl struct {
	// 句柄
	fd int
	// max event size, default: 100
	maxEventSize int
}

func (impl *epollImpl) Add(fd int) error {
	return unix.EpollCtl(impl.fd,
		unix.EPOLL_CTL_ADD,
		fd,
		&unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
}

func (impl *epollImpl) Remove(fd int) error {
	return unix.EpollCtl(impl.fd, syscall.EPOLL_CTL_DEL, fd, nil)
}

func (impl *epollImpl) Wait() ([]int, error) {
	events := make([]unix.EpollEvent, impl.maxEventSize)
	n, err := unix.EpollWait(impl.fd, events, 100)
	if err != nil {
		return nil, err
	}

	fds := make([]int, n)
	for i := 0; i < n; i++ {
		if events[i].Fd == 0 {
			continue
		}
		fds[i] = int(events[i].Fd)
	}
	return fds, nil
}

func Create() (*epollImpl, error) {
	fd, err := unix.EpollCreate(1)
	if err != nil {
		return nil, err
	}

	return &epollImpl{
		fd:           fd,
		maxEventSize: 100,
	}, nil
}
