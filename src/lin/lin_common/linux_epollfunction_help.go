//go:build linux
// +build linux

package lin_common

import "golang.org/x/sys/unix"

type EPOLL_EVENT int
const (
	EPOLL_EVENT_READ EPOLL_EVENT = 1
	EPOLL_EVENT_WRITE EPOLL_EVENT = 2
	EPOLL_EVENT_ALL EPOLL_EVENT = EPOLL_EVENT_READ | EPOLL_EVENT_WRITE
)

const (
	// EPOLLET:go not support epoll et mod
	_epoll_BASE_EVENTS = unix.EPOLLPRI | unix.EPOLLERR /*| unix.EPOLLHUP | unix.EPOLLRDHUP | unix.EPOLLET*/
	_epoll_READ_EVENTS =  unix.EPOLLIN
	_epoll_WRITE_EVENTS = unix.EPOLLOUT
)
func unixEpollCreate()(int, error) {
	return unix.EpollCreate1(unix.EPOLL_CLOEXEC)
}

func unixEpollAdd(efd int, fd int, evtInput EPOLL_EVENT, userData int32) error {
	var eEvent uint32 = _epoll_BASE_EVENTS
	if (evtInput & EPOLL_EVENT_READ) != 0 {
		eEvent |= _epoll_READ_EVENTS
	}
	if (evtInput & EPOLL_EVENT_WRITE) != 0 {
		eEvent |= _epoll_WRITE_EVENTS
	}
	evt := &unix.EpollEvent{Fd: int32(fd), Events: eEvent, Pad:userData}
	return unix.EpollCtl(efd, unix.EPOLL_CTL_ADD, fd, evt)
}

func unixEpollDel(efd int, fd int) error {
	return unix.EpollCtl(efd, unix.EPOLL_CTL_DEL, fd, nil)
}

func unixEpollMod(efd int, fd int, evtInput EPOLL_EVENT, userData int32) error {
	var eEvent uint32 = _epoll_BASE_EVENTS
	if (evtInput & EPOLL_EVENT_READ) != 0 {
		eEvent |= _epoll_READ_EVENTS
	}
	if (evtInput & EPOLL_EVENT_WRITE) != 0 {
		eEvent |= _epoll_WRITE_EVENTS
	}
	evt := &unix.EpollEvent{Fd: int32(fd), Events: eEvent, Pad:userData}
	return unix.EpollCtl(efd, unix.EPOLL_CTL_MOD, fd, evt)
}
