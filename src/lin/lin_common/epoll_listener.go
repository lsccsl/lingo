//go:build linux
// +build linux

package lin_common

import (
	"bytes"
	"golang.org/x/sys/unix"
	"math/rand"
	"runtime"
)

func ConstructorTcpConnectionInfo(fd int, magic int32, isDial bool, buffInitLen int)*TcpConnectionInfo {
	unix.SetNonblock(fd, true)
	ti := &TcpConnectionInfo{
		_fd: fd,
		_readBuf : bytes.NewBuffer(make([]byte, 0, buffInitLen)),
		_writeBuf : bytes.NewBuffer(make([]byte, 0, buffInitLen)),
		_magic: magic,
		_isDial: isDial,
		_cur_epoll_evt : EPOLL_EVENT_READ,
	}
	if isDial {
		ti._isConnSuc = false
		ti._cur_epoll_evt = EPOLL_EVENT_READ_WRITE
	} else {
		ti._isConnSuc = true
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

	unixEpollDel(pthis._epollFD, fd)
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
				magic := pthis._lsn.EPollListenerGenMagic()
				LogDebug("new conn fd:", t._fdConn, " magic:", magic)
				ti := ConstructorTcpConnectionInfo(t._fdConn, magic, false, pthis._lsn._paramTcpRWBuffLen)
				pthis._mapTcp[t._fdConn] = ti
				unixEpollAdd(pthis._epollFD, t._fdConn, ti._cur_epoll_evt, magic)

				if pthis._lsn._cb != nil {
					func(){
						defer func() {
							err := recover()
							if err != nil {
								LogErr(err)
							}
						}()

						pthis._lsn._cb.TcpAcceptConnection(ti._fd, ti._magic, ti._addr)
					}()
				}
			}

		case *Event_TcpClose:
			{
				LogDebug("recv user close tcp, fd:", t._fd, " magic:", t._magic)
				pthis.EpollConnection_close_tcp(t._fd, t._magic)
			}

		case *Event_TcpDial:
			{
				LogDebug("dial tcp connection, fd:", t._fd, " magic:", t._magic)
				ti := ConstructorTcpConnectionInfo(t._fd, t._magic, true, pthis._lsn._paramTcpRWBuffLen)
				pthis._mapTcp[t._fd] = ti
				unixEpollAdd(pthis._epollFD, t._fd, ti._cur_epoll_evt, t._magic) // if the tcp connection can write, means the tcp connection is success, it will be mod epoll wait read event when connection is ok
			}

		case *Event_TcpWrite:
			{
				LogDebug(" user tcp write, fd:", t._fdConn, " magic:", t._magic)
				pthis.EpollConnection_user_write(t._fdConn, t._magic, t._binData)
			}

		default:
		}
	}
}

func (pthis*EPollConnection)EpollConnection_user_write(fd int, magic int32, binData []byte) {
	ti, _ := pthis._mapTcp[fd]
	if ti == nil {
		return
	}

	if ti._magic != magic {
		LogDebug("magic not match, fd:", fd, " magic:", ti._magic, " old magic:", magic)
		return
	}

	ti._writeBuf.Write(binData)
	pthis.EpollConnection_do_write(ti, 100)
}

func (pthis*EPollConnection)EpollConnection_epllEvt_tcpread(fd int, magic int32, maxReadcount int) {
	ti, _ := pthis._mapTcp[fd]
	if ti == nil {
		return
	}
	if ti._magic != magic {
		LogDebug("magic not match, fd:", fd, " magic:", ti._magic, " old magic:", magic)
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

func (pthis*EPollConnection)EpollConnection_epllEvt_tcpwrite(fd int, magic int32){
	ti, _ := pthis._mapTcp[fd]
	if ti == nil {
		return
	}

	if ti._magic != magic {
		LogDebug("magic not match, fd:", fd, " magic:", ti._magic, " old magic:", magic)
		return
	}

	if ti._isDial && !ti._isConnSuc {
		ti._isConnSuc = true
		ti._cur_epoll_evt = EPOLL_EVENT_READ
		unixEpollMod(pthis._epollFD, fd, ti._cur_epoll_evt, magic)
		pthis._lsn._cb.TcpDialConnection(fd, ti._magic, ti._addr)
	}

	pthis.EpollConnection_do_write(ti, 100)
}

func (pthis*EPollConnection)EpollConnection_do_write(ti *TcpConnectionInfo, maxWriteLoop int) {
	// do write, if all data is write success, mod to epoll wait read
	if ti == nil {
		return
	}

	bModEpoll := true

	if ti._writeBuf.Len() != 0 {
		for i := 0; i < maxWriteLoop; i ++ {
			write_sz, err, bAgain := _tcpWrite(ti._fd, ti._writeBuf.Bytes())
			if err != nil {
				// write fail, will close tcp connection
				pthis.EpollConnection_close_tcp(ti._fd, ti._magic)
				return
			}
			ti._writeBuf.Next(write_sz)

			if ti._writeBuf.Len() == 0 {
				// all data has been write to buffer, change mod to epoll wait read
				break
			}

			if bAgain {
				// not all data write to buffer, write buffer is full, need change to epoll wait write mod
				ti._cur_epoll_evt = EPOLL_EVENT_READ_WRITE
				unixEpollMod(pthis._epollFD, ti._fd, ti._cur_epoll_evt, ti._magic)
				bModEpoll = false
				break
			}
		}
	}

	if bModEpoll {
		if ti._cur_epoll_evt != EPOLL_EVENT_READ {
			ti._cur_epoll_evt = EPOLL_EVENT_READ
			unixEpollMod(pthis._epollFD, ti._fd, ti._cur_epoll_evt, ti._magic)
		}
	}
}

func (pthis*EPollConnection)EPollConnection_AddEvent(evt interface{}) {
	pthis._evtQue.Enqueue(evt)
	unix.Write(pthis._evtFD, EVENT_BIN_1)
}

func (pthis*EPollConnection)_go_EpollConnection_epollwait() {
	defer func() {
		err := recover()
		if err != nil {
			LogErr(err)
		}
		LogErr("_go_EpollConnection_epollwait quit")
		pthis._lsn = nil
	}()

	events := make([]unix.EpollEvent, pthis._lsn._paramMaxEpollEventCount) // todo: change the events array size by epoll wait ret count
	for {
		count, err := unix.EpollWait(pthis._epollFD, events, pthis._lsn._paramEpollWaitTimeoutMills)
		if count == 0 || (count < 0 && err == unix.EINTR) {
			runtime.Gosched()
			continue
		} else if err != nil {
			LogErr("epoll wait err:", err)
			break
		}

		for i := 0; i < count; i ++ {
			epollEvt := &events[i]
			triggerFD := int(epollEvt.Fd)
			if triggerFD == pthis._evtFD {
				pthis.EpollConnection_process_evt()
			} else {
				// tcp read / write / err
				if (epollEvt.Events & unix.EPOLLIN) != 0 {
					pthis.EpollConnection_epllEvt_tcpread(triggerFD, epollEvt.Pad, 100)
				}
				if (epollEvt.Events & unix.EPOLLOUT) != 0 {
					pthis.EpollConnection_epllEvt_tcpwrite(triggerFD, epollEvt.Pad)
				}
				if (epollEvt.Events & unix.EPOLLERR) != 0 {
					pthis.EpollConnection_close_tcp(triggerFD, epollEvt.Pad)
				}
			}
		}
	}
}


func (pthis*EPollAccept)_go_EpollAccept_epollwait() {
	defer func() {
		err := recover()
		if err != nil {
			LogErr(err)
		}
		pthis._lsn = nil
		LogDebug("_go_EpollAccept_epollwait quit")
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
		err = unixEpollAdd(pthis._epollAccept._epollFD, pthis._epollAccept._tcpListenerFD, EPOLL_EVENT_READ, 0)
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
		err = unixEpollAdd(pthis._epollAccept._epollFD, pthis._epollAccept._evtFD, EPOLL_EVENT_READ, 0)
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
			err = unixEpollAdd(epollConn._epollFD, epollConn._evtFD, EPOLL_EVENT_READ, 0)
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

func (pthis*EPollListener)EPollListenerGenMagic() int32 {
	return int32(rand.Int()*rand.Int())
}

func (pthis*EPollListener)EPollListenerAddEvent(fd int, evt interface{}) {
	if len(pthis._epollConnection) == 0 {
		LogErr("epoll connection coroutine count is 0")
		return
	}
	idx := fd % len(pthis._epollConnection)
	if idx >= len(pthis._epollConnection) {
		return
	}
	pthis._epollConnection[idx].EPollConnection_AddEvent(evt)
}

func (pthis*EPollListener)EPollListenerCloseTcp(rawfd int, magic int32){
	pthis.EPollListenerAddEvent(rawfd, &Event_TcpClose{rawfd, magic})
}

func (pthis*EPollListener)EPollListenerWrite(rawfd int, magic int32, binData []byte) {
	pthis.EPollListenerAddEvent(rawfd, &Event_TcpWrite{rawfd, magic, binData})
}

func (pthis*EPollListener)EPollListenerDial(addr string)(rawfd int, magic int32, err error){
	LogDebug(" begin connect addr:", addr)
	rawfd, err = _tcpConnectNoBlock(addr)
	if err != nil {
		return -1, 0, err
	}

	magic = pthis.EPollListenerGenMagic()
	pthis.EPollListenerAddEvent(rawfd,&Event_TcpDial{rawfd,magic})

	LogDebug(" connect fd:", rawfd, " magic:", magic, " addr:", addr)
	return rawfd, magic, nil
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
