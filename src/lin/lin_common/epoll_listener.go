//go:build linux
// +build linux

package lin_common

import (
	"golang.org/x/sys/unix"
	"sync"
	"unsafe"
)

const MTU = 1536

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

	_binRead []byte
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
	_readBufLen int

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
				unixEpollAdd(pthis._epollFD, t._fdConn, epoll_READ_EVENTS, 0)
			}
		default:
		}
	}
}


func (pthis*EPollConnection)EpollConnection_tcpread(fd int, maxReadcount int) {
	for i := 0; i < maxReadcount; i ++{
		n, err := _tcpRead(fd, pthis._binRead)

		if err != nil {
			LogDebug("tcp read err:", err, " fd:", fd)
			// close tcp
			break
		}

		if n == 0 { // no data
			break
		}

		// read until no data any more
		LogDebug("tcpread:", n)
	}

	/*
		from c++
		do
		{
			if(pfi->rpos_ >= pfi->rbuf_.size())
			{
				pfi->rbuf_.resize(pfi->rbuf_.size() + 1024);
			}

			int32 ret = CChannel::TcpRead(pfi->fd_, &pfi->rbuf_[pfi->rpos_], pfi->rbuf_.size() - pfi->rpos_);

			MYLOG_DEBUG(("read ret:%d", ret));

			if(ret < 0)
			{
				MYLOG_DEBUG(("err need close"));
				need_close = 1;
				break;
			}
			else if(0 == ret)
			{
				MYLOG_DEBUG(("no data"));
				break;
			}
			else
			{
				pfi->rpos_ += ret;

				pfi->byte_read_ += ret;

				MYLOG_DUMP_BIN(&pfi->rbuf_[0], pfi->rpos_);

			}
		}while(1);
	*/
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
				pthis.EpollConnection_tcpread(triggerFD, 100)
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
	maxEpollEventCount int, epollWaitTimeoutMills int, readBufLen int) (*EPollListener, error){
	if epollCoroutineCount <= 0 {
		epollCoroutineCount = 1
	}

	el := &EPollListener{
		_maxEpollEventCount : maxEpollEventCount,
		_epollWaitTimeoutMills : epollWaitTimeoutMills,
		_readBufLen: readBufLen,
	}
	el._epollAccept._lsn = el
	el._epollAccept._evtQue = NewLKQueue()
	el._epollAccept._evtBuf = make([]byte, EVENT_BIN_8)

	if el._readBufLen <= 0 {
		el._readBufLen = MTU
	}

	var err error

	{
		// create epoll fd
		el._epollAccept._epollFD, err = unixEpollCreate()
		if err != nil {
			return nil, GenErrNoERR_NUM("create epoll accept handle fail:", err)
		}
		// create tcp listener fd
		el._epollAccept._tcpListenerFD, err = _tcpListen(addr)
		if err != nil {
			return nil, err
		}

		// add tcp listener fd to epoll wait
		err = unixEpollAdd(el._epollAccept._epollFD, el._epollAccept._tcpListenerFD, epoll_READ_EVENTS, 0)
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
		err = unixEpollAdd(el._epollAccept._epollFD, el._epollAccept._evtFD, epoll_READ_EVENTS, 0)
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
			_binRead : make([]byte, el._readBufLen),
		}
		el._epollConnection = append(el._epollConnection, epollConn)
		epollConn._epollFD, err = unixEpollCreate()
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
			err = unixEpollAdd(epollConn._epollFD, epollConn._evtFD, epoll_READ_EVENTS, 0)
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