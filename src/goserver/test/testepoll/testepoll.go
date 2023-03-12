//go:build linux
// +build linux

package main

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/unix"
	"goserver/common"
	"net"
	"time"
)

var ipv4InIPv6Prefix = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}

func sockaddrInet4ToIP(sa *unix.SockaddrInet4) net.IP {
	ip := make([]byte, 16)
	// ipv4InIPv6Prefix
	copy(ip[0:12], ipv4InIPv6Prefix)
	copy(ip[12:16], sa.Addr[:])
	return ip
}
func sockaddrInet6ToIPAndZone(sa *unix.SockaddrInet6) (net.IP, string) {
	ip := make([]byte, 16)
	copy(ip, sa.Addr[:])
	return ip, ip6ZoneToString(int(sa.ZoneId))
}
func ip6ZoneToString(zone int) string {
	if zone == 0 {
		return ""
	}
	if ifi, err := net.InterfaceByIndex(zone); err == nil {
		return ifi.Name
	}
	return int2decimal(uint(zone))
}

// Convert int to decimal string.
func int2decimal(i uint) string {
	if i == 0 {
		return "0"
	}

	// Assemble decimal in reverse order.
	b := make([]byte, 32)
	bp := len(b)
	for ; i > 0; i /= 10 {
		bp--
		b[bp] = byte(i%10) + '0'
	}
	return string(b[bp:])
}

func SockaddrToTCPOrUnixAddr(sa unix.Sockaddr) net.Addr {
	switch sa := sa.(type) {
	case *unix.SockaddrInet4:
		ip := sockaddrInet4ToIP(sa)
		return &net.TCPAddr{IP: ip, Port: sa.Port}
	case *unix.SockaddrInet6:
		ip, zone := sockaddrInet6ToIPAndZone(sa)
		return &net.TCPAddr{IP: ip, Port: sa.Port, Zone: zone}
	case *unix.SockaddrUnix:
		return &net.UnixAddr{Name: sa.Name, Net: "unix"}
	}
	return nil
}

func tcpListen(addr string) int {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	fmt.Println(tcpAddr, err)
	fmt.Println("ip:", tcpAddr.IP, " port:", &tcpAddr.Port)

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|/*unix.SOCK_NONBLOCK|*/unix.SOCK_CLOEXEC, unix.IPPROTO_TCP)
	fmt.Println(fd, err)

	sa4 := &unix.SockaddrInet4{Port: tcpAddr.Port}
	if tcpAddr.IP != nil {
		if len(tcpAddr.IP) == 16 {
			copy(sa4.Addr[:], tcpAddr.IP[12:16]) // copy last 4 bytes of slice to array
		} else {
			copy(sa4.Addr[:], tcpAddr.IP) // copy all bytes of slice to array
		}
	}
	err = unix.Bind(fd, sa4)
	fmt.Println("bind:", err)

	err = unix.Listen(fd, 128)
	fmt.Println("listen:", err)

	return fd
}

func tcpAccept(listenFD int) (connFD int, addr string) {
	connFD, sa, err := unix.Accept(listenFD)
	fmt.Println("new connect:", connFD, sa, err)

	remoteAddr := SockaddrToTCPOrUnixAddr(sa)
	fmt.Println(remoteAddr)

	addr = remoteAddr.String()
	return
}

func epollAddRead(efd int, sockfd int){
	evt := &unix.EpollEvent{Fd: int32(sockfd), Events: unix.EPOLLPRI | unix.EPOLLIN}
	err := unix.EpollCtl(efd, unix.EPOLL_CTL_ADD, sockfd, evt)
	fmt.Println("epoll ctl:", err)
}

func epollWait(efd int) {
	events := make([]unix.EpollEvent, 128)
	count, err := unix.EpollWait(efd, events, 30000)
	fmt.Println("epoll wait:", count, err)
}

func main() {
	common.InitLog("testepoll.log", "testepoll_err.log", true, false, false)
	testepoll()
	efd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	fmt.Println(efd, " err:", err)

	fd := tcpListen("192.168.2.129:3001")
	epollAddRead(efd, fd)

	for {
		epollWait(efd)
		connFD, addr := tcpAccept(fd)
		fmt.Println("new connection:", connFD, addr)
	}

	for {
		time.Sleep(time.Second * 10)
	}
}

type test_tcp_info struct {
	fd common.FD_DEF
}
type test_cb struct {
	lsn *common.EPollListener
	mapFD map[int]*test_tcp_info
}
func (pthis*test_cb)TcpAcceptConnection(fd common.FD_DEF, addr net.Addr, attachData interface{})interface{}{
	common.LogDebug("new tcp connection:", fd.FD, " addr:", addr)
	ti := &test_tcp_info{fd: fd}
	pthis.mapFD[fd.FD]=ti
	return nil
}
func (pthis*test_cb)TcpDialConnection(fd common.FD_DEF, addr net.Addr, attachData interface{})interface{}{
	common.LogDebug("suc to dial, fd:", fd.FD, " magic:", fd.Magic, " addr:", addr)
	return nil
}

func (pthis*test_cb)TcpData(fd common.FD_DEF, readBuf *bytes.Buffer, attachData interface{})(bytesProcess int, retAttachData interface{}) {
	common.LogDebug("tcp data:", fd.FD, " data len:", readBuf.Len())
	ti := pthis.mapFD[fd.FD]
	if ti != nil {
		if fd.FD%2 != 0 {
			common.LogDebug("will write back fd:", fd.FD, " magic:", fd.Magic)
			//pthis.lsn.EPollListenerWrite(fd, readBuf.Bytes())
			//pthis.lsn.EPollListenerCloseTcp(rawfd, ti.magic)
		}
	}
	return readBuf.Len(), nil
}

func (pthis*test_cb)TcpClose(fd common.FD_DEF, reason common.EN_TCP_CLOSE_REASON, attachData interface{}){
	common.LogDebug("tcp close fd:", fd.FD)
}
func (pthis*test_cb)TcpOutBandData(fd common.FD_DEF, data interface{}, inAttachData interface{}) {

}

func (pthis*test_cb)TcpTick(fd common.FD_DEF, tNowMill int64, inAttachData interface{}){
	common.LogDebug(fd, " tNowMill:", tNowMill, " inAttachData:", inAttachData)
}
func testepoll() {
	tcb := &test_cb{
		mapFD : make(map[int]*test_tcp_info),
	}
	el, err := common.ConstructorEPollListener(tcb,"192.168.2.129:3001", 10, common.ParamEPollListener{})
	tcb.lsn = el
	fmt.Println("lin_common.ConstructEPollListener", el, err)
	fd, err := tcb.lsn.EPollListenerDial("192.168.2.129:2003", nil, false)
	common.LogDebug("fd:", fd.FD, " magic:", fd.Magic, " err:", err)
	el.EPollListenerWait()
}