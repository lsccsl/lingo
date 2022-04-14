package lin_common

import (
	"golang.org/x/sys/unix"
	"sync"
)


type EPollConnectionCoroutine struct {
	_epollFD int

	_lsn *EPollListener
}
type EPollAcceptCoroutine struct {
	_epollFD int
	_tcpListenerFD int

	_lsn *EPollListener
}

type EPollListener struct {
	EpollAccept EPollAcceptCoroutine
	EpollConnection []*EPollConnectionCoroutine

	_maxEpollEventCount int
	_epollWaitTimeoutMills int

	_wg sync.WaitGroup
}


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
			// tcp read or write
		}
	}
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

	var err error
	// create epoll fd
	el.EpollAccept._epollFD, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		return nil, GenErrNoERR_NUM("create epoll accept handle fail")
	}
	// create tcp listener fd
	el.EpollAccept._tcpListenerFD, err = _tcpListen(addr)
	if err != nil {
		return nil, err
	}
	evt := &unix.EpollEvent{Fd: int32(el.EpollAccept._tcpListenerFD), Events: unix.EPOLLPRI | unix.EPOLLIN}
	err = unix.EpollCtl(el.EpollAccept._epollFD, unix.EPOLL_CTL_ADD, el.EpollAccept._tcpListenerFD, evt)
	if err != nil {
		return nil, GenErrNoERR_NUM("add listener fd to epoll fail")
	}
	el._wg.Add(1)
	go el.EpollAccept._goEpollAcceptCoroutine()

	for i := 0; i < epollCoroutineCount; i ++ {
		epollConn := &EPollConnectionCoroutine{_lsn: el}
		epollConn._epollFD, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
		if err != nil {
			return nil, GenErrNoERR_NUM("create epoll connection handle fail")
		}
		el._wg.Add(1)
		go epollConn._goEpollConnectionCoroutine()
	}

	return el, nil
}

func (pthis*EPollListener)EPollListenerWait() {
	pthis._wg.Wait()
}