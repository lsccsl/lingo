package lin_common

import "golang.org/x/sys/unix"

type EPOLL_EVENT int
const (
	EPOLL_EVENT_READ EPOLL_EVENT = 1
	EPOLL_EVENT_WRITE EPOLL_EVENT = 2
	EPOLL_EVENT_ALL EPOLL_EVENT = EPOLL_EVENT_READ | EPOLL_EVENT_WRITE
)

const (
	epoll_READ_EVENTS = unix.EPOLLPRI | unix.EPOLLIN
	epoll_WRITEE_VENTS = unix.EPOLLOUT
	epoll_READWRITE_EVENTS = epoll_READ_EVENTS | epoll_WRITEE_VENTS
)
func unixEpollCreate()(int, error) {
	return unix.EpollCreate1(unix.EPOLL_CLOEXEC)
}

func unixEpollAdd(efd int, fd int, evtInput EPOLL_EVENT, userData int32) error {
	var eEvent uint32 = 0
	if (evtInput & EPOLL_EVENT_READ) != 0 {
		eEvent |= epoll_READ_EVENTS
	}
	if (evtInput & EPOLL_EVENT_WRITE) != 0 {
		eEvent |= epoll_WRITEE_VENTS
	}
	evt := &unix.EpollEvent{Fd: int32(fd), Events: eEvent, Pad:userData}
	return unix.EpollCtl(efd, unix.EPOLL_CTL_ADD, fd, evt)
}

func unixEpollDel(efd int, fd int) error {
	return unix.EpollCtl(efd, unix.EPOLL_CTL_DEL, fd, nil)
}
