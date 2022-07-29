//go:build linux
// +build linux

package lin_common

import (
	"bytes"
	"golang.org/x/sys/unix"
	"math/rand"
	"runtime"
)

func ConstructorTcpConnectionInfo(fd FD_DEF, isDial bool, buffInitLen int, needKeepAlive bool, attachData interface{})*tcpConnectionInfo {
	err := _setNoBlock(fd.FD)
	if err != nil {
		LogDebug("_setNoBlock:", fd.String(), " err:", err)
	}
	err = _setLingerOff(fd.FD)
	if err != nil {
		LogDebug("_setLingerOff:", fd.String(), " err:", err)
	}
	err = _setNoDelay(fd.FD)
	if err != nil {
		LogDebug("_setNoDelay:", fd.String(), " err:", err)
	}
	if needKeepAlive {
		err = _tcpKeepAlive(fd.FD, 10, 10, 10)
		if err != nil {
			LogDebug("_setNoDelay:", fd.String(), " err:", err)
		}
	}
	err = _setRecvBuffer(fd.FD, 65535)
	if err != nil {
		LogErr("_setRecvBuffer:", fd.String(), " err:", err)
	}
	err = _setSendBuffer(fd.FD, 65535)
	if err != nil {
		LogErr("_setSendBuffer:", fd.String(), " err:", err)
	}

	ti := &tcpConnectionInfo{
		_fd: fd,
		_readBuf : bytes.NewBuffer(make([]byte, 0, buffInitLen)),
		_writeBuf : bytes.NewBuffer(make([]byte, 0, buffInitLen)),
		_isDial: isDial,
		_cur_epoll_evt : EPOLL_EVENT_READ,
		_attachData: attachData,
		_closeReason: EN_TCP_CLOSE_REASON_inter_none,
	}
	if isDial {
		ti._isConnSuc = false
		ti._cur_epoll_evt = EPOLL_EVENT_READ_WRITE
	} else {
		ti._isConnSuc = true
	}
	ti._addr = _tcpGetPeerName(ti._fd.FD)

	return ti
}

func (pthis*ePollConnection)_add_tcp_conn(ti*tcpConnectionInfo) {
	pthis._mapTcp[ti._fd.FD] = ti

	pthis._tcpConnCount = len(pthis._mapTcp)
}
func (pthis*ePollConnection)_del_tcp_conn(fd int) {
	delete(pthis._mapTcp, fd)

	pthis._tcpConnCount = len(pthis._mapTcp)
}
func (pthis*ePollConnection)_get_tcp_conn(fd int)*tcpConnectionInfo {
	ti, ok := pthis._mapTcp[fd]
	if ti == nil || !ok {
		return nil
	}
	return ti
}

func (pthis*ePollConnection)EpollConnection_close_tcp(fd FD_DEF, reason EN_TCP_CLOSE_REASON){
	ti := pthis._get_tcp_conn(fd.FD)
	if ti == nil{
		//LogDebug(" can't find tcp:", fd.String(), " ti:", ti)
		return
	}

	if ti._fd.Magic != fd.Magic {
		// linux fd will auto rollback, new tcp connection fd will take the slot whole that last fd closed
		LogDebug("magic not match fd:", fd, " current:", ti._fd.String(), " close:", fd.String())
		return
	}

	unixEpollDel(pthis._epollFD, fd.FD)
	if pthis._lsn._cb != nil {
		func(){
			defer func() {
				err := recover()
				if err != nil {
					LogErr(err)
				}
			}()
			pthis._lsn._cb.TcpClose(ti._fd, reason, ti._attachData)
		}()
	}

	pthis._del_tcp_conn(fd.FD)
	unix.Close(fd.FD)
	pthis._tcpCloseCount ++
}

func (pthis*ePollConnection)EpollConnection_process_evt(){
	for {
		evt := pthis._evtQue.Dequeue()
		if evt == nil {
			break
		}
		unix.Read(pthis._evtFD, pthis._evtBuf)

		switch t:=evt.(type){
		case *event_NewConnection:
			{
				// new tcp connection add to epoll
				magic := pthis._lsn.EPollListenerGenMagic()
				//LogDebug("new conn fd:", t._fdConn, " magic:", magic)
				ti := ConstructorTcpConnectionInfo(FD_DEF{t._fdConn, magic}, false, pthis._lsn._paramTcpRWBuffLen, false, nil)
				pthis._add_tcp_conn(ti)
				unixEpollAdd(pthis._epollFD, t._fdConn, ti._cur_epoll_evt, magic, pthis._lsn._paramET)

				if pthis._lsn._cb != nil {
					func(){
						defer func() {
							err := recover()
							if err != nil {
								LogErr(err)
							}
						}()

						ad := pthis._lsn._cb.TcpAcceptConnection(ti._fd, ti._addr, ti._attachData)
						if ad != nil {
							ti._attachData = ad
						}
					}()
				}
			}

		case *event_TcpClose:
			{
				LogDebug("recv user close tcp, fd:", t.fd.FD, " magic:", t.fd.Magic)
				pthis.EpollConnection_close_tcp(t.fd, t.reason)
			}

		case *event_TcpDial:
			{
				//LogDebug("dial tcp connection, fd:", t.fd.FD, " magic:", t.fd.Magic)
				ti := ConstructorTcpConnectionInfo(t.fd, true, pthis._lsn._paramTcpRWBuffLen, false, t.attachData)
				pthis._add_tcp_conn(ti)
				unixEpollAdd(pthis._epollFD, t.fd.FD, ti._cur_epoll_evt, t.fd.Magic, pthis._lsn._paramET) // if the tcp connection can write, means the tcp connection is success, it will be mod epoll wait read event when connection is ok
			}

		case *event_TcpWrite:
			{
				//LogDebug(" user tcp write, fd:", t.fd.FD, " magic:", t.fd.Magic)
				pthis.EpollConnection_user_write(t.fd, t._binData)
			}

		default:
		}
	}

/*	pthis._evt_need_process_next_loop = false
	atomic.StoreInt64(&pthis._evt_process, 0)
	if ((!pthis._evtQue.IsEmpty()) && atomic.CompareAndSwapInt64(&pthis._evt_process, 0, 1)) {
		_, err := unix.Write(pthis._evtFD, EVENT_BIN_1)
		if err != unix.EAGAIN {
			pthis._evt_need_process_next_loop = true
		}
	}*/
}

func (pthis*ePollConnection)EpollConnection_user_write(fd FD_DEF, binData []byte) {
	ti := pthis._get_tcp_conn(fd.FD)
	if ti == nil {
		return
	}

	if ti._fd.Magic != fd.Magic {
		LogDebug("magic not match, fd:", fd.String(), " current fd:", ti._fd.String())
		return
	}

	ti._writeBuf.Write(binData)
	pthis.EpollConnection_do_write(ti)
}

func (pthis*ePollConnection)EpollConnection_epllEvt_tcpread(fd FD_DEF) {
	ti := pthis._get_tcp_conn(fd.FD)
	if ti == nil {
		return
	}
	if ti._fd.Magic != fd.Magic {
		LogDebug("magic not match, fd:", fd.FD, " magic:", ti._fd.Magic, " old magic:", fd.Magic)
		return
	}

	maxReadcount := pthis._lsn._paramMaxTcpRead
	if pthis._lsn._paramET {
		maxReadcount = -1
	} else {
		if maxReadcount <= 0 {
			maxReadcount = 1
		}
	}

	//LogDebug("fd:", ti._fd.String(), " maxReadcount:", maxReadcount)

	bClose := false

	// if not support epoll et mode, can set read count limited, if support epoll et mode, maxReadcount must big enought
	for i := 0; (i < maxReadcount) || (maxReadcount < 0); i ++{
		n, err := _tcpRead(fd.FD, pthis._binRead)

		if err != nil {
			LogDebug("tcp read err fd:", fd.String(), " err:", err)
			// close tcp, del from epoll
			bClose = true
			break
		}

		if n == 0 { // no data
			break
		}

		pthis._byteRecv += int64(n)

		// read until no data anymore
		ti._readBuf.Write(pthis._binRead[:n])
		//LogDebug("tcp read fd:", fd, " count:", n, " read buf len:", ti._readBuf.Len())
	}

	if ti._readBuf.Len() > 0 {
		if pthis._lsn._cb != nil {
			func(){
				defer func() {
					err := recover()
					if err != nil {
						LogErr(err)
					}
				}()
				for ti._readBuf.Len() > 0 {
					bytesProcess, attachData := pthis._lsn._cb.TcpData(ti._fd, ti._readBuf, ti._attachData)
					if bytesProcess <= 0 {
						break
					}

					if attachData != nil {
						ti._attachData = attachData
					}

					pthis._byteProc += int64(bytesProcess)

					ti._readBuf.Next(bytesProcess)
				}
			}()
		} else {
			LogDebug("no call back fd:", ti._fd.FD)
			ti._readBuf.Next(ti._readBuf.Len())
		}
	}

	if ti._readBuf.Len() == 0 {
		ti._readBuf.Reset() // todo : shrink buf
	}

	if bClose {
		pthis.EpollConnection_close_tcp(fd, EN_TCP_CLOSE_REASON_read_err)
		return
	}
}

func (pthis*ePollConnection)EpollConnection_epllEvt_tcpwrite(fd FD_DEF){
	ti := pthis._get_tcp_conn(fd.FD)
	if ti == nil {
		return
	}

	if ti._fd.Magic != fd.Magic {
		LogDebug("magic not match, fd:", fd.FD, " magic:", ti._fd.Magic, " old magic:", fd.Magic)
		return
	}

	if ti._isDial && !ti._isConnSuc {
		ti._isConnSuc = true
		ti._cur_epoll_evt = EPOLL_EVENT_READ
		unixEpollMod(pthis._epollFD, ti._fd.FD, ti._cur_epoll_evt, ti._fd.Magic, pthis._lsn._paramET)
		ti._addr = _tcpGetPeerName(ti._fd.FD)
		ad := pthis._lsn._cb.TcpDialConnection(ti._fd, ti._addr, ti._attachData)
		if ad != nil {
			ti._attachData = ad
		}
	}

	pthis.EpollConnection_do_write(ti)
}

func (pthis*ePollConnection)EpollConnection_do_write(ti *tcpConnectionInfo) {
	// do write, if all data is write success, mod to epoll wait read
	if ti == nil {
		return
	}

	bModEpoll := true

	if ti._writeBuf.Len() != 0 {

		maxWriteLoop := pthis._lsn._paramMaxTcpWrite
		if pthis._lsn._paramET {
			maxWriteLoop = -1
		} else {
			if maxWriteLoop <= 0 {
				maxWriteLoop = 1
			}
		}

		//LogDebug("fd:", ti._fd.String(), " maxWriteLoop:", maxWriteLoop)

		// if not support epoll et mode, can set read count limited, if support epoll et mode, maxReadcount must big enought
		for i := 0; (i < maxWriteLoop) || (maxWriteLoop < 0); i ++ {
			write_sz, err, bAgain := _tcpWrite(ti._fd.FD, ti._writeBuf.Bytes())
			LogDebug("write:", write_sz, " err:", err, " bAgain:", bAgain)
			if err != nil {
				// write fail, will close tcp connection
				pthis.EpollConnection_close_tcp(ti._fd, EN_TCP_CLOSE_REASON_write_err)
				return
			}

			pthis._byteSend += int64(write_sz)

			if write_sz > 0 {
				ti._writeBuf.Next(write_sz)
			}

			if ti._writeBuf.Len() == 0 {
				// all data has been write to buffer, change mod to epoll wait read
				LogDebug("write done")
				break
			}

			if bAgain {
				// not all data write to buffer, write buffer is full, need change to epoll wait write mod
				ti._cur_epoll_evt = EPOLL_EVENT_READ_WRITE
				unixEpollMod(pthis._epollFD, ti._fd.FD, ti._cur_epoll_evt, ti._fd.Magic, pthis._lsn._paramET)
				bModEpoll = false
				break
			}
		}
	}

	if ti._writeBuf.Len() == 0 {
		ti._writeBuf.Reset() // todo : shrink buf ti._writeBuf = bytes.NewBuffer(make([]byte, ti._writeBuf.Cap()/2))
	}

	if bModEpoll {
		if ti._cur_epoll_evt != EPOLL_EVENT_READ {
			ti._cur_epoll_evt = EPOLL_EVENT_READ
			unixEpollMod(pthis._epollFD, ti._fd.FD, ti._cur_epoll_evt, ti._fd.Magic, pthis._lsn._paramET)
		}
	}
}

func (pthis*ePollConnection)EPollConnection_AddEvent(fd int, evt interface{}) {
	pthis._evtQue.Enqueue(evt)
	n, err := unix.Write(pthis._evtFD, EVENT_BIN_1)
	if err != nil || n < len(EVENT_BIN_1) {
		LogErr("fd:", fd, " write event err:", err, " n:", n)
	}
/*	if atomic.CompareAndSwapInt64(&pthis._evt_process, 0, 1) {
		unix.Write(pthis._evtFD, EVENT_BIN_1)
	}*/
}

func (pthis*ePollConnection)_go_EpollConnection_epollwait() {
	defer func() {
		err := recover()
		if err != nil {
			LogErr(err)
		}
		pthis._lsn = nil
		LogErr("_go_EpollConnection_epollwait quit")
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
			fd := FD_DEF{triggerFD, epollEvt.Pad}
			if triggerFD == pthis._evtFD {
				pthis.EpollConnection_process_evt()
			} else {
				// tcp read / write / err
				if (epollEvt.Events & unix.EPOLLIN) != 0 {
					pthis.EpollConnection_epllEvt_tcpread(fd)
				}
				if (epollEvt.Events & unix.EPOLLOUT) != 0 {
					pthis.EpollConnection_epllEvt_tcpwrite(fd)
				}
				if (epollEvt.Events & unix.EPOLLERR) != 0 {
					pthis.EpollConnection_close_tcp(fd, EN_TCP_CLOSE_REASON_epoll_err)
				}
			}
		}

/*		if pthis._evt_need_process_next_loop {
			pthis.EpollConnection_process_evt()
		}*/
	}
}


func (pthis*ePollAccept)_go_EpollAccept_epollwait() {
	defer func() {
		err := recover()
		if err != nil {
			LogErr(err)
		}
		pthis._lsn = nil
		LogErr("_go_EpollAccept_epollwait quit")
	}()

	maxReadcount := pthis._lsn._paramMaxTcpRead
	if pthis._lsn._paramET {
		maxReadcount = -1
	} else {
		if maxReadcount <= 0 {
			maxReadcount = 1
		}
	}

	LogDebug("accept maxReadcount:", maxReadcount)

	events := make([]unix.EpollEvent, pthis._lsn._paramMaxEpollEventCount)
	for {
		count, err := unix.EpollWait(pthis._epollFD, events, pthis._lsn._paramEpollWaitTimeoutMills) // todo:改成select,此处不需要用epoll
		if count == 0 || (count < 0 && err == unix.EINTR) {
			runtime.Gosched()
			continue
		} else if err != nil {
			LogErr("epoll wait err:", err)
			break
		}

		for i := 0; i < count && i < len(events); i ++ {
			triggerFD := int(events[i].Fd)
			if triggerFD == pthis._evtFD {
				continue
			}
			// tcp accept

			for i:= 0; (i < maxReadcount) || (maxReadcount < 0); i ++ {
				fd, addr, err := _tcpAccept(triggerFD)
				if fd < 0 {
					if err != nil{
						LogDebug("fail accept:", err, " fd:", fd, " listen fd:", triggerFD)
					}
					break
				}

				//LogDebug("new tcp connection, fd:", fd, " addr:", addr)
				pthis._lsn._EPollListenerAddEvent(fd, &event_NewConnection{_fdConn: fd, addr:addr})
			}
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

	if pthis._paramTmpReadBufLen <= 0 {
		pthis._paramTmpReadBufLen = MTU
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
		err = unixEpollAdd(pthis._epollAccept._epollFD, pthis._epollAccept._tcpListenerFD, EPOLL_EVENT_READ, 0, pthis._paramET)
		if err != nil {
			return GenErrNoERR_NUM("add listener fd to epoll fail:", err)
		}
	}

	{
		// create event fd
		pthis._epollAccept._evtFD, err = unixEvent()
		if err != nil {
			return err
		}

		// add event fd to epoll wait
		err = unixEpollAdd(pthis._epollAccept._epollFD, pthis._epollAccept._evtFD, EPOLL_EVENT_READ, 0, pthis._paramET)
		if err != nil {
			return GenErrNoERR_NUM("add listener fd to epoll fail:", err)
		}
	}

	pthis._wg.Add(1)
	go pthis._epollAccept._go_EpollAccept_epollwait()

	for i := 0; i < epollCoroutineCount; i ++ {
		epollConn := &ePollConnection{
			_lsn: pthis,
			_evtQue:NewLKQueue(),
			_evtBuf : make([]byte, EVENT_BIN_8),
			_binRead : make([]byte, pthis._paramTmpReadBufLen),
			_mapTcp : make(MAP_TCPCONNECTION),
		}
		pthis._epollConnection = append(pthis._epollConnection, epollConn)
		epollConn._epollFD, err = unixEpollCreate()
		if err != nil {
			return GenErrNoERR_NUM("create epoll connection handle fail:", err)
		}

		{
			// create event fd
			epollConn._evtFD, err = unixEvent()
			if err != nil {
				return err
			}

			// add event fd to epoll wait
			err = unixEpollAdd(epollConn._epollFD, epollConn._evtFD, EPOLL_EVENT_READ, 0, pthis._paramET)
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

func (pthis*EPollListener)_EPollListenerAddEvent(fd int, evt interface{}) {
	connCount := len(pthis._epollConnection)
	if connCount <= 0 {
		LogErr("epoll connection coroutine count is 0")
		return
	}
	idx := fd % connCount
	if idx >= connCount {
		return
	}
	pthis._epollConnection[idx].EPollConnection_AddEvent(fd, evt)
}

func (pthis*EPollListener)EPollListenerCloseTcp(fd FD_DEF, reason EN_TCP_CLOSE_REASON){
	pthis._EPollListenerAddEvent(fd.FD, &event_TcpClose{fd:fd, reason:reason})
}

func (pthis*EPollListener)EPollListenerWrite(fd FD_DEF, binData []byte) {
	LogDebug("write data len:", len(binData))
	pthis._EPollListenerAddEvent(fd.FD, &event_TcpWrite{fd:fd, _binData:binData})
}

func (pthis*EPollListener)EPollListenerDial(addr string, attachData interface{})(fd FD_DEF, err error){
	//LogDebug(" begin connect addr:", addr)
	rawfd, err := _tcpConnectNoBlock(addr)
	if err != nil {
		LogErr(" fail connect addr:", addr)
		return FD_DEF{-1,0}, err
	}

	magic := pthis.EPollListenerGenMagic()
	pthis._EPollListenerAddEvent(rawfd,&event_TcpDial{
		fd: FD_DEF{rawfd,magic},
		attachData: attachData,
	})

	//LogDebug(" connect fd:", rawfd, " magic:", magic, " addr:", addr)
	return FD_DEF{rawfd, magic}, nil
}

func(pthis*EPollListener)EPollListenerGetStatic(es *EPollListenerStatic) {
	for _, val := range pthis._epollConnection {
		es.TcpConnCount += val._tcpConnCount
		es.TcpCloseCount += val._tcpCloseCount
		es.ByteRecv += val._byteRecv
		es.ByteProc += val._byteProc
		es.ByteSend += val._byteSend
	}
}
