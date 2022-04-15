//go:build linux
// +build linux

package lin_common

import (
	"golang.org/x/sys/unix"
	"sync"
	"unsafe"
)

const (
	EPOLL_READ_EVENTS = unix.EPOLLPRI | unix.EPOLLIN
	EPOLL_WRITEE_VENTS = unix.EPOLLOUT
	EPOLL_READWRITE_EVENTS = EPOLL_READ_EVENTS | EPOLL_WRITEE_VENTS
)

var (
	EVENT_1 uint64 = 1
	EVENT_BIN_1 = (*(*[8]byte)(unsafe.Pointer(&EVENT_1)))[:]
	EVENT_BIN_8 = 8
)

type EPollEventNewConnection struct {
	_fdConn int
}

type EPollConnection struct {
	_epollFD int

	_evtFD int
	_evtQue *LKQueue // bind for _evtFD todo:改成用go自带的锁队列
	_evtBuf []byte

	_lsn *EPollListener
}
type EPollAccept struct {
	_epollFD int // todo:改成select
	_tcpListenerFD int

	_evtFD int
	_evtQue *LKQueue // bind for _evtFD todo:改成用go自带的锁队列
	_evtBuf []byte

	_lsn *EPollListener
}

type EPollListener struct {
	_epollAccept EPollAccept
	_epollConnection []*EPollConnection

	_maxEpollEventCount int
	_epollWaitTimeoutMills int

	_wg sync.WaitGroup
}

func (pthis*EPollConnection)EpollConnection_process_evt(){
	unix.Read(pthis._evtFD, pthis._evtBuf)
	for {
		evt := pthis._evtQue.Dequeue()
		if evt == nil {
			break
		}

		switch t:=evt.(type){
		case *EPollEventNewConnection:
			{
				// new tcp connection add to epoll
				LogDebug("new conn:", t)
				evt := &unix.EpollEvent{Fd: int32(t._fdConn), Events: EPOLL_READ_EVENTS}
				unix.EpollCtl(pthis._epollFD, unix.EPOLL_CTL_ADD, t._fdConn, evt)
			}
		default:
		}
	}
}

func (pthis*EPollConnection)_goEpollConnectionCoroutine() {
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
				pthis.EpollConnection_process_evt()
			} else {
				// tcp read or write
			}
		}
	}
}

func (pthis*EPollConnection)EPollConnection_AddEvent(evt interface{}) {
	pthis._evtQue.Enqueue(evt)
	unix.Write(pthis._evtFD, EVENT_BIN_1)
}

func (pthis*EPollAccept)_goEpollAcceptCoroutine() {
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
			pthis._lsn.EPollListenerAddEvent(fd, &EPollEventNewConnection{_fdConn: fd})
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
	el._epollAccept._lsn = el
	el._epollAccept._evtQue = NewLKQueue()
	el._epollAccept._evtBuf = make([]byte, EVENT_BIN_8)

	var err error

	{
		// create epoll fd
		el._epollAccept._epollFD, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
		if err != nil {
			return nil, GenErrNoERR_NUM("create epoll accept handle fail:", err)
		}
		// create tcp listener fd
		el._epollAccept._tcpListenerFD, err = _tcpListen(addr)
		if err != nil {
			return nil, err
		}

		// add tcp listener fd to epoll wait
		evt := &unix.EpollEvent{Fd: int32(el._epollAccept._tcpListenerFD), Events: EPOLL_READ_EVENTS}
		err = unix.EpollCtl(el._epollAccept._epollFD, unix.EPOLL_CTL_ADD, el._epollAccept._tcpListenerFD, evt)
		if err != nil {
			return nil, GenErrNoERR_NUM("add listener fd to epoll fail:", err)
		}
	}

	{
		// create event fd
		el._epollAccept._evtFD, err = _linuxEvent()
		if err != nil {
			return nil, err
		}

		// add event fd to epoll wait
		evt := &unix.EpollEvent{Fd: int32(el._epollAccept._evtFD), Events: EPOLL_READ_EVENTS}
		err = unix.EpollCtl(el._epollAccept._epollFD, unix.EPOLL_CTL_ADD, el._epollAccept._evtFD, evt)
		if err != nil {
			return nil, GenErrNoERR_NUM("add listener fd to epoll fail:", err)
		}
	}

	el._wg.Add(1)
	go el._epollAccept._goEpollAcceptCoroutine()

	for i := 0; i < epollCoroutineCount; i ++ {
		epollConn := &EPollConnection{
			_lsn: el,
			_evtQue:NewLKQueue(),
			_evtBuf : make([]byte, EVENT_BIN_8),
		}
		el._epollConnection = append(el._epollConnection, epollConn)
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

func (pthis*EPollListener)EPollListenerAddEvent(fd int, evt interface{}) {
	idx := fd % len(pthis._epollConnection)
	if idx >= len(pthis._epollConnection) {
		return
	}
	pthis._epollConnection[idx].EPollConnection_AddEvent(evt)
}