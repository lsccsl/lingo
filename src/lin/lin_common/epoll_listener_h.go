//go:build linux
// +build linux

package lin_common

import (
	"bytes"
	"net"
	"strconv"
	"sync"
	"unsafe"
)

const MTU = 1536

var (
	EVENT_1 uint64 = 1
	EVENT_BIN_1 = (*(*[8]byte)(unsafe.Pointer(&EVENT_1)))[:]
	EVENT_BIN_8 = 8
)


type FD_DEF struct {
	FD int
	Magic int32
}
func (pthis*FD_DEF)String()string{
	return "fd:" + strconv.Itoa(pthis.FD) + " magic:" + strconv.Itoa(int(pthis.Magic))
}
func (pthis*FD_DEF)IsSame(fd *FD_DEF) bool {
	if pthis.FD != fd.FD {
		return false
	}
	if pthis.Magic != fd.Magic {
		return false
	}
	return true
}

/* @brief begin inter evetn define */
type event_NewConnection struct { // new tcp connection event
	_fdConn int
}
type event_TcpWrite struct { // tcp write event
	fd FD_DEF
	_binData []byte
}
type event_TcpClose struct { // user close tcp connection
	fd FD_DEF
}
type event_TcpDial struct {
	fd FD_DEF
}
/* @brief end inter evetn define */


type EPollCallback interface {
	TcpAcceptConnection(fd FD_DEF, addr net.Addr)
	TcpDialConnection(fd FD_DEF, addr net.Addr)
	TcpData(fd FD_DEF, readBuf *bytes.Buffer)(bytesProcess int)
	TcpClose(fd FD_DEF)
}


/* @brief tcp connection info define */
type tcpConnectionInfo struct {
	_readBuf *bytes.Buffer
	_writeBuf *bytes.Buffer
	_fd FD_DEF
	_addr net.Addr

	_cur_epoll_evt EPOLL_EVENT

	_isDial bool
	_isConnSuc bool
}
type MAP_TCPCONNECTION map[int]*tcpConnectionInfo


type ePollConnection_Interface interface {
	EpollConnection_process_evt()
	EpollConnection_epllEvt_tcpread(fd FD_DEF)
	EpollConnection_epllEvt_tcpwrite(fd FD_DEF)
	EpollConnection_user_write(fd FD_DEF, binData []byte)
	EpollConnection_do_write(ti *tcpConnectionInfo)
	EPollConnection_AddEvent(evt interface{})
	EpollConnection_close_tcp(fd FD_DEF)
	_go_EpollConnection_epollwait()

	_add_tcp_conn(*tcpConnectionInfo)
	_del_tcp_conn(fd int)
	_get_tcp_conn(fd int)*tcpConnectionInfo
}
type ePollConnectionStatic struct {
	_tcpConnCount int
	_tcpCloseCount int64
	_byteRecv int64
	_byteProc int64
	_byteSend int64
}
type ePollConnection struct {
	_epollFD int

	_evtFD int
	_evtQue *LKQueue // bind for _evtFD todo:改成用go自带的锁队列
	_evtBuf []byte
/*	_evt_process int64
	_evt_need_process_next_loop bool*/

	_lsn *EPollListener

	_binRead []byte
	_mapTcp MAP_TCPCONNECTION

	ePollConnectionStatic
}


type ePollAccept_interface interface {
	_go_EpollAccept_epollwait()
}
type ePollAccept struct {
	_epollFD int // todo:改成select
	_tcpListenerFD int

	_evtFD int
	_evtQue *LKQueue // bind for _evtFD todo:改成用go自带的锁队列
	_evtBuf []byte

	_lsn *EPollListener
}


type ParamEPollListener struct {
	ParamMaxEpollEventCount int
	ParamEpollWaitTimeoutMills int
	ParamTmpReadBufLen int    // epoll coroutine tmp read buf
	ParamTcpRWBuffLen  int // tcp r/w data buffer
	ParamMaxTcpRead int
	ParamMaxTcpWrite int
	ParamET bool
}
type interParamEPollListener struct {
	_paramMaxEpollEventCount int
	_paramEpollWaitTimeoutMills int
	_paramTmpReadBufLen int    // epoll coroutine tmp read buf
	_paramTcpRWBuffLen  int // tcp r/w data buffer
	_paramMaxTcpRead int
	_paramMaxTcpWrite int
	_paramET bool // if support epoll et mode
}

type EPollListenerStatic struct {
	TcpConnCount int
	TcpCloseCount int64
	ByteRecv int64
	ByteProc int64
	ByteSend int64
}

type EPollListener_interface interface {
	EPollListenerInit(cb EPollCallback, addr string, epollCoroutineCount int) error
	EPollListenerWait()
	EPollListenerAddEvent(fd int, evt interface{})
	EPollListenerCloseTcp(fd FD_DEF)
	EPollListenerWrite(fd FD_DEF, binData []byte)
	EPollListenerDial(addr string)(fd FD_DEF, err error)
	EPollListenerGetStatic(*EPollListenerStatic)
}

// EPollListener : epoll application interface
type EPollListener struct {
	_epollAccept ePollAccept
	_epollConnection []*ePollConnection

	_cb EPollCallback
	interParamEPollListener

	_wg sync.WaitGroup
}
func ConstructorEPollListener(cb EPollCallback, addr string, epollCoroutineCount int,
	param ParamEPollListener) (*EPollListener, error){
	el := &EPollListener{
		//ParamEPollListener:param,
		interParamEPollListener:interParamEPollListener{
			_paramMaxEpollEventCount : param.ParamMaxEpollEventCount,
			_paramEpollWaitTimeoutMills : param.ParamEpollWaitTimeoutMills,
			_paramTmpReadBufLen : param.ParamTmpReadBufLen,
			_paramTcpRWBuffLen : param.ParamTcpRWBuffLen,
			_paramMaxTcpRead : param.ParamMaxTcpRead,
			_paramMaxTcpWrite : param.ParamMaxTcpWrite,
			_paramET : param.ParamET,
		},
		_cb : cb,
	}

	if el._paramET {
		el._paramMaxTcpRead = -1
		el._paramMaxTcpWrite = -1
	}

	if el._paramMaxEpollEventCount <= 0 {
		el._paramMaxEpollEventCount = 2048
	}
	if el._paramEpollWaitTimeoutMills <= 0 {
		el._paramEpollWaitTimeoutMills = 300 * 1000
	}
	if el._paramTmpReadBufLen <= 0 {
		el._paramTmpReadBufLen = MTU
	}
	if el._paramTcpRWBuffLen <=0 {
		el._paramTcpRWBuffLen = 8192
	}
	if el._paramMaxTcpRead <= 0 {
		el._paramMaxTcpRead = 100
	}
	if el._paramMaxTcpWrite <= 0 {
		el._paramMaxTcpWrite = 100
	}

	return el, el.EPollListenerInit(cb, addr, epollCoroutineCount)
}
