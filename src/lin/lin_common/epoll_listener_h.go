package lin_common

import (
	"bytes"
	"net"
	"sync"
	"unsafe"
)

const MTU = 1536

var (
	EVENT_1 uint64 = 1
	EVENT_BIN_1 = (*(*[8]byte)(unsafe.Pointer(&EVENT_1)))[:]
	EVENT_BIN_8 = 8
)


/* @brief begin inter evetn define */
type Event_NewConnection struct { // new tcp connection event
	_fdConn int
}
type Event_TcpWrite struct { // tcp write event
	_fdConn int
	_binData []byte
}
type Event_TcpClose struct { // user close tcp connection
	_fd int
	_magic int32
}
type Event_TcpDial struct {
	_fd int
	_magic int32
}
/* @brief end inter evetn define */


type EPollCallback interface {
	TcpAcceptConnection(rawfd int, magic int32, addr net.Addr)
	TcpDialConnection(rawfd int, magic int32, addr net.Addr)
	TcpData(rawfd int, readBuf *bytes.Buffer)(bytesProcess int)
	TcpClose(rawfd int)
}


/* @brief tcp connection info define */
type TcpConnectionInfo struct {
	_readBuf *bytes.Buffer
	_writeBuf *bytes.Buffer
	_fd int
	_addr net.Addr

	_magic int32

	_isDial bool
	_isConnSuc bool
}
type MAP_TCPCONNECTION map[int]*TcpConnectionInfo


type EPollConnection_Interface interface {
	EpollConnection_process_evt()
	EpollConnection_tcpread(fd int, magic int32, maxReadcount int)
	EPollConnection_AddEvent(evt interface{})
	EpollConnection_close_tcp(fd int, magic int32)
	_go_EpollConnection_epollwait()
}
type EPollConnection struct {
	_epollFD int

	_evtFD int
	_evtQue *LKQueue // bind for _evtFD todo:改成用go自带的锁队列
	_evtBuf []byte

	_lsn *EPollListener

	_binRead []byte
	_mapTcp MAP_TCPCONNECTION
}


type EPollAccept_interface interface {
	_go_EpollAccept_epollwait()
}
type EPollAccept struct {
	_epollFD int // todo:改成select
	_tcpListenerFD int

	_evtFD int
	_evtQue *LKQueue // bind for _evtFD todo:改成用go自带的锁队列
	_evtBuf []byte

	_lsn *EPollListener
}


type EPollListener_interface interface {
	EPollListenerInit(cb EPollCallback, addr string, epollCoroutineCount int) error
	EPollListenerWait()
	EPollListenerAddEvent(fd int, evt interface{})
	EPollListenerCloseTcp(rawfd int, magic int32)
	EPollListenerDial(addr string)(rawfd int, magic int32, err error)
}
type EPollListener struct {
	_epollAccept EPollAccept
	_epollConnection []*EPollConnection

	_cb EPollCallback

	_paramMaxEpollEventCount int
	_paramEpollWaitTimeoutMills int
	_paramReadBufLen int
	_paramTcpRWBuffLen  int


	_wg sync.WaitGroup
}
