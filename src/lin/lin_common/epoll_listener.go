//go:build linux
// +build linux

package lin_common

import (
	"golang.org/x/sys/unix"
	"sync"
	"unsafe"
)


type EPollConnectionCoroutine struct {
	_epollFD int
	_evtFD int
	_evtQue *LKQueue // bind for _evtFD todo:改成用go自带的锁队列

	_lsn *EPollListener
}
type EPollAcceptCoroutine struct {
	_epollFD int
	_tcpListenerFD int
	_evtFD int
	_evtQue *LKQueue // bind for _evtFD todo:改成用go自带的锁队列

	_lsn *EPollListener
}

type EPollListener struct {
	EpollAccept EPollAcceptCoroutine
	EpollConnection []*EPollConnectionCoroutine

	_maxEpollEventCount int
	_epollWaitTimeoutMills int

	_wg sync.WaitGroup
}

const (
	EPOLL_READ_EVENTS = unix.EPOLLPRI | unix.EPOLLIN
	EPOLL_WRITEE_VENTS = unix.EPOLLOUT
	EPOLL_READWRITE_EVENTS = EPOLL_READ_EVENTS | EPOLL_WRITEE_VENTS
)

var (
	EVENT_1 uint64 = 1
	EVENT_BIN_1 = (*(*[8]byte)(unsafe.Pointer(&EVENT_1)))[:]
)

func (pthis*EPollConnectionCoroutine)_goEpollConnectionCoroutine() {
	defer func() {
		pthis._lsn = nil
		err := recover()
		if err != nil {
			LogErr(err)
		}
	}()


	events := make([]unix.EpollEvent, pthis._lsn._maxEpollEventCount) // todo: change the events array size by epoll wait ret count
	for {
		count, err := unix.EpollWait(pthis._epollFD, events, pthis._lsn._epollWaitTimeoutMills)
		if err != nil {
			LogErr("epoll wait err")
			break
		}

		for i := 0; i < count; i ++ {
			triggerFD := int(events[i].Fd)
			if triggerFD == pthis._evtFD {
				// read from chan
			} else {
				// tcp read or write
			}
		}
	}
}

func (pthis*EPollConnectionCoroutine)_EPollConnectionCoroutineAddEvent(evt interface{}) {
	pthis._evtQue.Enqueue(evt)
	unix.Write(pthis._evtFD, EVENT_BIN_1)
}

func (pthis*EPollAcceptCoroutine)_goEpollAcceptCoroutine() {
	defer func() {
		pthis._lsn = nil
		err := recover()
		if err != nil {
			LogErr(err)
		}
	}()

	events := make([]unix.EpollEvent, pthis._lsn._maxEpollEventCount) // todo: change the events array size by epoll wait ret count
	for {
		count, err := unix.EpollWait(pthis._epollFD, events, pthis._lsn._epollWaitTimeoutMills)
		if err != nil {
			LogErr("epoll wait err")
			break
		}

		for i := 0; i < count && i < len(events); i ++ {
			triggerFD := int(events[i].Fd)
			if triggerFD == pthis._evtFD {
				continue
			}
			// tcp accept
			fd, addr, err := _tcpAccept(int(events[i].Fd))
			if err != nil {
				LogErr("fail accept")
				continue
			}

			LogDebug("fd:", fd, " addr:", addr)
			//add connection epoll wait coroutine by addr hash
		}
	}
}

func ConstructEPollListener(addr string, epollCoroutineCount int,
	maxEpollEventCount int, epollWaitTimeoutMills int) (*EPollListener, error){
	if epollCoroutineCount <= 0 {
		epollCoroutineCount = 1
	}

	el := &EPollListener{
		_maxEpollEventCount : maxEpollEventCount,
		_epollWaitTimeoutMills : epollWaitTimeoutMills,
	}
	el.EpollAccept._lsn = el
	el.EpollAccept._evtQue = NewLKQueue()

	var err error

	{
		// create epoll fd
		el.EpollAccept._epollFD, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
		if err != nil {
			return nil, GenErrNoERR_NUM("create epoll accept handle fail:", err)
		}
		// create tcp listener fd
		el.EpollAccept._tcpListenerFD, err = _tcpListen(addr)
		if err != nil {
			return nil, err
		}

		// add tcp listener fd to epoll wait
		evt := &unix.EpollEvent{Fd: int32(el.EpollAccept._tcpListenerFD), Events: EPOLL_READ_EVENTS}
		err = unix.EpollCtl(el.EpollAccept._epollFD, unix.EPOLL_CTL_ADD, el.EpollAccept._tcpListenerFD, evt)
		if err != nil {
			return nil, GenErrNoERR_NUM("add listener fd to epoll fail:", err)
		}
	}

	{
		// create event fd
		el.EpollAccept._evtFD, err = _linuxEvent()
		if err != nil {
			return nil, err
		}

		// add event fd to epoll wait
		evt := &unix.EpollEvent{Fd: int32(el.EpollAccept._evtFD), Events: EPOLL_READ_EVENTS}
		err = unix.EpollCtl(el.EpollAccept._epollFD, unix.EPOLL_CTL_ADD, el.EpollAccept._evtFD, evt)
		if err != nil {
			return nil, GenErrNoERR_NUM("add listener fd to epoll fail:", err)
		}
	}

	el._wg.Add(1)
	go el.EpollAccept._goEpollAcceptCoroutine()

	for i := 0; i < epollCoroutineCount; i ++ {
		epollConn := &EPollConnectionCoroutine{
			_lsn: el,
			_evtQue:NewLKQueue(),
		}
		epollConn._epollFD, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
		if err != nil {
			return nil, GenErrNoERR_NUM("create epoll connection handle fail:", err)
		}

		{
			// create event fd
			epollConn._evtFD, err = _linuxEvent()
			if err != nil {
				return nil, err
			}

			// add event fd to epoll wait
			evt := &unix.EpollEvent{Fd: int32(epollConn._evtFD), Events: EPOLL_READ_EVENTS}
			err = unix.EpollCtl(epollConn._epollFD, unix.EPOLL_CTL_ADD, epollConn._evtFD, evt)
			if err != nil {
				return nil, GenErrNoERR_NUM("add listener fd to epoll fail:", err)
			}
		}
		el._wg.Add(1)
		go epollConn._goEpollConnectionCoroutine()
	}

	return el, nil
}

func (pthis*EPollListener)EPollListenerWait() {
	pthis._wg.Wait()
}