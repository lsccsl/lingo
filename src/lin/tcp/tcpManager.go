package tcp

import (
	"fmt"
	"lin/lin_common"
	"net"
	"strconv"
	"sync"
)

type TcpMgr struct {
	tcpLsn       net.Listener
	wg             sync.WaitGroup
	cbConnection   InterfaceTcpConnection
	closeExpireSec int

	TcpDialMgr

	mapConnMutex sync.Mutex
	mapConn      MAP_TCPCONN
}

func (pthis *TcpMgr) CBGenConnectionID() TCP_CONNECTION_ID {
	return TCP_CONNECTION_ID(lin_common.GenUUID64_V4())
}
func (pthis *TcpMgr) CBAddTcpConn(tcpConn *TcpConnection) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	pthis.mapConn[tcpConn.TcpConnectionID()] = tcpConn
}
func (pthis *TcpMgr) CBGetConnectionCB() InterfaceTcpConnection {
	return pthis.cbConnection
}
func (pthis *TcpMgr) CBDelTcpConn(id TCP_CONNECTION_ID) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	delete(pthis.mapConn, id)
}

func (pthis *TcpMgr)GetTcpConnection(tcpConnID TCP_CONNECTION_ID) *TcpConnection {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()
	conn, _ := pthis.mapConn[tcpConnID]
	return conn
}

func (pthis *TcpMgr)go_tcpAccept() {
	for {
		conn, err := pthis.tcpLsn.Accept()
		if err != nil {
			lin_common.LogErr("tcp accept err:", err)
			if conn != nil {
				conn.Close()
			}
			continue
		}

		if conn == nil {
			lin_common.LogErr(" tcp conn is nil")
			continue
		}

		_, err = startTcpConnection(pthis, conn, pthis.closeExpireSec)
		if err != nil {
			lin_common.LogErr("start accept tcp connect err", err)
		}
	}

	pthis.wg.Done()
}

func StartTcpManager(ip string, port int, CBConnection InterfaceTcpConnection,  closeExpireSec int) (*TcpMgr, error) {
	t := &TcpMgr{}

	addr := ip + ":" + strconv.Itoa(port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	t.tcpLsn = lsn
	t.cbConnection = CBConnection
	t.closeExpireSec = closeExpireSec
	t.mapConn = make(MAP_TCPCONN)

	t.wg.Add(1)
	go t.go_tcpAccept()

	t.TcpDialMgrStart(t, closeExpireSec)

	return t, nil
}

func (pthis *TcpMgr) TcpMgrWait() {
	lin_common.LogDebug("begin wait")
	pthis.wg.Wait()
	lin_common.LogDebug("end wait")
}


func (pthis *TcpMgr) TcpMgrCloseConn(id TCP_CONNECTION_ID) {
	conn := pthis.GetTcpConnection(id)
	if conn == nil {
		return
	}
	conn.TcpConnectClose()
}

func (pthis *TcpMgr) TcpMgrDump(bDtail bool) (str string, totalRecv int64, totalSend int64, totalProc int64){
	str = "\r\ntcp connect:\r\n"
	func(){
		pthis.mapConnMutex.Lock()
		defer pthis.mapConnMutex.Unlock()
		mapUnprocessd := make(map[TCP_CONNECTION_ID]int)
		for _, val := range pthis.mapConn {
			if bDtail {
				var remoteAddr string
				var localAddr string
				if val.netConn != nil {
					remoteAddr = val.netConn.RemoteAddr().String()
				}
				if val.netConn != nil {
					localAddr = val.netConn.LocalAddr().String()
				}
				str += fmt.Sprintf(" \r\n connection:%v remote:[%v] local:[%v] IsAccept:%v srv:%v ClientID:%v"+
					" recv:%v send:%v proc:%v SendCount:%v",
					val.TcpConnectionID(), remoteAddr, localAddr, val.IsAccept, val.SrvID, val.ClientID,
					val.ByteRecv, val.ByteSend, val.ByteProc, val.SendCount)
			}
			totalRecv += val.ByteRecv
			totalProc += val.ByteProc
			totalSend += val.ByteSend
			if val.ByteRecv != val.ByteProc {
				mapUnprocessd[val.TcpConnectionID()] = int(val.ByteRecv - val.ByteProc)
			}
		}
		str += "\r\ntcp conn count:" + strconv.Itoa(len(pthis.mapConn))
		//str += fmt.Sprintf(" not process bytes:%v\r\n", mapUnprocessd)
		str += fmt.Sprintf(" not process client:%v totalRecv:%v totalProc:%v totalSend:%v unprocess:%v",
			len(mapUnprocessd), totalRecv, totalProc, totalSend, totalRecv - totalProc)
	}()

	str += "\r\ntcp dial data\r\n"
	func(){
		pthis.TcpDialMgr.mapDialDataMutex.Lock()
		defer pthis.TcpDialMgr.mapDialDataMutex.Unlock()
		for _, val := range pthis.TcpDialMgr.mapDialData {
			if bDtail {
				str += fmt.Sprintf("\r\n srv:%v connection:%v [%v:%v]", val.srvID, val.tcpConnID, val.ip, val.port)
			}
		}
		str += fmt.Sprintf("\r\n dial count:%v \r\n", len(pthis.TcpDialMgr.mapDialData))
	}()

	return
}
