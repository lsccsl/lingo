//go:build linux
// +build linux

package lin_common

import (
	"bytes"
	"golang.org/x/sys/unix"
	"math/rand"
)

func ConstructorTcpConnectionInfo(fd int, magic int32, buffInitLen int)*TcpConnectionInfo {
	unix.SetNonblock(fd, true)
	ti := &TcpConnectionInfo{
		_fd: fd,
		_readBuf : bytes.NewBuffer(make([]byte, 0, buffInitLen)),
		_writeBuf : bytes.NewBuffer(make([]byte, 0, buffInitLen)),
		_magic: magic,
	}
	sa, err := unix.Getpeername(ti._fd)
	if err == nil {
		ti._addr = _sockaddrToTCPOrUnixAddr(sa)
	}

	return ti
}

func (pthis*EPollConnection)EpollConnection_close_tcp(fd int, magic int32){
	ti, ok := pthis._mapTcp[fd]
	if ti == nil || !ok{
		return
	}

	if ti._magic != magic {
		// linux fd will auto rollback, new tcp connection fd will take the slot whole that last fd closed
		LogDebug("magic not match fd:", fd, " current:", ti._magic, " close:", magic)
		return
	}

	if pthis._lsn._cb != nil {
		func(){
			defer func() {
				err := recover()
				if err != nil {
					LogErr(err)
				}
			}()
			pthis._lsn._cb.TcpClose(fd)
		}()
	}

	unixEpollDel(pthis._epollFD, fd)
	delete(pthis._mapTcp, fd)

	unix.Close(fd)
}

func (pthis*EPollConnection)EpollConnection_process_evt(){
	unix.Read(pthis._evtFD, pthis._evtBuf)
	for {
		evt := pthis._evtQue.Dequeue()
		if evt == nil {
			break
		}

		switch t:=evt.(type){
		case *Event_NewConnection:
			{
				// new tcp connection add to epoll
				magic := int32(rand.Int() + rand.Int())
				LogDebug("new conn fd:", t._fdConn, " magic:", magic)
				unixEpollAdd(pthis._epollFD, t._fdConn, epoll_READ_EVENTS, magic)
				ti := ConstructorTcpConnectionInfo(t._fdConn, magic, pthis._lsn._paramTcpRWBuffLen)
				pthis._mapTcp[t._fdConn] = ti

				if pthis._lsn._cb != nil {
					func(){
						defer func() {
							err := recover()
							if err != nil {
								LogErr(err)
							}
						}()

						pthis._lsn._cb.TcpNewConnection(ti._fd, ti._magic, ti._addr)
					}()
				}
			}

		case *Event_TcpClose:
			{
				LogDebug("recv close tcp, fd:", t._fd, " magic:", t._magic)
				pthis.EpollConnection_close_tcp(t._fd, t._magic)
			}

		default:
		}
	}
}

func (pthis*EPollConnection)EpollConnection_tcpread(fd int, magic int32, maxReadcount int) {
	ti, _ := pthis._mapTcp[fd]
	if ti == nil {
		return
	}

	if maxReadcount <= 0 {
		maxReadcount = 1
	}

	bClose := false

	// not support epoll et mode, can set read count limited
	for i := 0; i < maxReadcount; i ++{
		n, err := _tcpRead(fd, pthis._binRead)

		if err != nil {
			LogDebug("tcp read err fd:", fd, " err:", err)
			// close tcp, del from epoll
			bClose = true
			break
		}

		if n == 0 { // no data
			break
		}

		// read until no data anymore
		ti._readBuf.Write(pthis._binRead[:n])
		LogDebug("tcp read fd:", fd, " count:", n, " read buf len:", ti._readBuf.Len())
	}

	if bClose {
		pthis.EpollConnection_close_tcp(fd, magic)
		return
	}

	if pthis._lsn._cb != nil {
		func(){
			defer func() {
				err := recover()
				if err != nil {
					LogErr(err)
				}
			}()
			for ti._readBuf.Len() > 0 {
				bytesProcess := pthis._lsn._cb.TcpData(ti._fd, ti._readBuf)
				if bytesProcess <= 0 {
					break
				}
				ti._readBuf.Next(bytesProcess)
			}
		}()
	} else {
		LogDebug("no call back fd:", ti._fd)
		ti._readBuf.Next(ti._readBuf.Len())
	}
}

func (pthis*EPollConnection)EPollConnection_AddEvent(evt interface{}) {
	pthis._evtQue.Enqueue(evt)
	LogDebug("add event:", evt)
	unix.Write(pthis._evtFD, EVENT_BIN_1)
}

func (pthis*EPollConnection)_go_EpollConnection_epollwait() {
	defer func() {
		pthis._lsn = nil
		err := recover()
		if err != nil {
			LogErr(err)
		}
	}()

	events := make([]unix.EpollEvent, pthis._lsn._paramMaxEpollEventCount) // todo: change the events array size by epoll wait ret count
	for {
		count, err := unix.EpollWait(pthis._epollFD, events, pthis._lsn._paramEpollWaitTimeoutMills)
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
				pthis.EpollConnection_tcpread(triggerFD, events[i].Pad, 100)
			}
		}
	}
}


func (pthis*EPollAccept)_go_EpollAccept_epollwait() {
	defer func() {
		pthis._lsn = nil
		err := recover()
		if err != nil {
			LogErr(err)
		}
	}()

	events := make([]unix.EpollEvent, pthis._lsn._paramMaxEpollEventCount) // todo: change the events array size by epoll wait ret count
	for {
		count, err := unix.EpollWait(pthis._epollFD, events, pthis._lsn._paramEpollWaitTimeoutMills)
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

			LogDebug("new tcp connection, fd:", fd, " addr:", addr)
			pthis._lsn.EPollListenerAddEvent(fd, &Event_NewConnection{_fdConn: fd})
		}
	}
}

func (pthis*EPollListener)EPollListenerInit(cb EPollCallback, addr string, epollCoroutineCount int) error{
	if cb == nil {
		return GenErrNoERR_NUM("EPollCallback is nil")
	}

	if epollCoroutineCount <= 0 {
		epollCoroutineCount = 1
	}

	pthis._epollAccept._lsn = pthis
	pthis._epollAccept._evtQue = NewLKQueue()
	pthis._epollAccept._evtBuf = make([]byte, EVENT_BIN_8)

	if pthis._paramReadBufLen <= 0 {
		pthis._paramReadBufLen = MTU
	}

	var err error

	{
		// create epoll fd
		pthis._epollAccept._epollFD, err = unixEpollCreate()
		if err != nil {
			return GenErrNoERR_NUM("create epoll accept handle fail:", err)
		}
		// create tcp listener fd
		pthis._epollAccept._tcpListenerFD, err = _tcpListen(addr)
		if err != nil {
			return err
		}

		// add tcp listener fd to epoll wait
		err = unixEpollAdd(pthis._epollAccept._epollFD, pthis._epollAccept._tcpListenerFD, epoll_READ_EVENTS, 0)
		if err != nil {
			return GenErrNoERR_NUM("add listener fd to epoll fail:", err)
		}
	}

	{
		// create event fd
		pthis._epollAccept._evtFD, err = uinuxEvent()
		if err != nil {
			return err
		}

		// add event fd to epoll wait
		err = unixEpollAdd(pthis._epollAccept._epollFD, pthis._epollAccept._evtFD, epoll_READ_EVENTS, 0)
		if err != nil {
			return GenErrNoERR_NUM("add listener fd to epoll fail:", err)
		}
	}

	pthis._wg.Add(1)
	go pthis._epollAccept._go_EpollAccept_epollwait()

	for i := 0; i < epollCoroutineCount; i ++ {
		epollConn := &EPollConnection{
			_lsn: pthis,
			_evtQue:NewLKQueue(),
			_evtBuf : make([]byte, EVENT_BIN_8),
			_binRead : make([]byte, pthis._paramReadBufLen),
			_mapTcp : make(MAP_TCPCONNECTION),
		}
		pthis._epollConnection = append(pthis._epollConnection, epollConn)
		epollConn._epollFD, err = unixEpollCreate()
		if err != nil {
			return GenErrNoERR_NUM("create epoll connection handle fail:", err)
		}

		{
			// create event fd
			epollConn._evtFD, err = uinuxEvent()
			if err != nil {
				return err
			}

			// add event fd to epoll wait
			err = unixEpollAdd(epollConn._epollFD, epollConn._evtFD, epoll_READ_EVENTS, 0)
			if err != nil {
				return GenErrNoERR_NUM("add listener fd to epoll fail:", err)
			}
		}
		pthis._wg.Add(1)
		go epollConn._go_EpollConnection_epollwait()
	}

	return nil
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

func (pthis*EPollListener)EPollListenerCloseTcp(rawfd int, magic int32){
	pthis.EPollListenerAddEvent(rawfd, &Event_TcpClose{rawfd, magic})
}

func ConstructorEPollListener(cb EPollCallback, addr string, epollCoroutineCount int,
	maxEpollEventCount int, epollWaitTimeoutMills int, readBufLen int, tcpRWBuffLen int) (*EPollListener, error){
	el := &EPollListener{
		_paramMaxEpollEventCount : maxEpollEventCount,
		_paramEpollWaitTimeoutMills : epollWaitTimeoutMills,
		_paramReadBufLen : readBufLen,
		_paramTcpRWBuffLen : tcpRWBuffLen,
		_cb : cb,
	}
	return el, el.EPollListenerInit(cb, addr, epollCoroutineCount)
}
