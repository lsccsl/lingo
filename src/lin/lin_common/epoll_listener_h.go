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
/* @brief end inter evetn define */


type EPollCallback interface {
	TcpNewConnection(rawfd int, addr net.Addr)
	TcpData(rawfd int, readBuf *bytes.Buffer)(bytesProcess int)
}


/* @brief tcp connection info define */
type TcpConnectionInfo struct {
	_readBuf *bytes.Buffer
	_writeBuf *bytes.Buffer
	_fd int
	_addr net.Addr
}
type MAP_TCPCONNECTION map[int]*TcpConnectionInfo


type EPollConnection_Interface interface {
	EpollConnection_process_evt()
	EpollConnection_tcpread(fd int, maxReadcount int)
	EPollConnection_AddEvent(evt interface{})
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
