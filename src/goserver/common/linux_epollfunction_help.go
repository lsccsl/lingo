//go:build linux
// +build linux

package common

import "golang.org/x/sys/unix"

type EPOLL_EVENT int
const (
	EPOLL_EVENT_READ       EPOLL_EVENT = 1
	EPOLL_EVENT_WRITE      EPOLL_EVENT = 2
	EPOLL_EVENT_READ_WRITE EPOLL_EVENT = EPOLL_EVENT_READ | EPOLL_EVENT_WRITE
)

const (
	_epoll_BASE_EVENTS = unix.EPOLLPRI | unix.EPOLLERR
	_epoll_BASE_EVENTS_ET = unix.EPOLLPRI | unix.EPOLLERR | unix.EPOLLHUP | unix.EPOLLRDHUP | unix.EPOLLET
	_epoll_READ_EVENTS =  unix.EPOLLIN
	_epoll_WRITE_EVENTS = unix.EPOLLOUT
)
func unixEpollCreate()(int, error) {
	return unix.EpollCreate(1)
	//return unix.EpollCreate1(unix.EPOLL_CLOEXEC)
}

func unixEpollAdd(efd int, fd int, evtInput EPOLL_EVENT, userData int32, bET bool) error {
	var eEvent uint32 = _epoll_BASE_EVENTS
	if bET {
		//LogDebug("fd:", fd, " et mode")
		eEvent = _epoll_BASE_EVENTS_ET
	}
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

func unixEpollMod(efd int, fd int, evtInput EPOLL_EVENT, userData int32, bET bool) error {
	var eEvent uint32 = _epoll_BASE_EVENTS
	if bET {
		//LogDebug("fd:", fd, " et mode")
		eEvent = _epoll_BASE_EVENTS_ET
	}

	if (evtInput & EPOLL_EVENT_READ) != 0 {
		eEvent |= _epoll_READ_EVENTS
	}
	if (evtInput & EPOLL_EVENT_WRITE) != 0 {
		eEvent |= _epoll_WRITE_EVENTS
	}
	evt := &unix.EpollEvent{Fd: int32(fd), Events: eEvent, Pad:userData}
	return unix.EpollCtl(efd, unix.EPOLL_CTL_MOD, fd, evt)
}
