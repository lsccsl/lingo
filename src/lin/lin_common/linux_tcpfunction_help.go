//go:build linux
// +build linux

package lin_common

import (
	"golang.org/x/sys/unix"
	"net"
	"syscall"
)

var _ipv4InIPv6Prefix = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}

func _sockaddrInet4ToIP(sa *unix.SockaddrInet4) net.IP {
	ip := make([]byte, 16)
	// ipv4InIPv6Prefix
	copy(ip[0:12], _ipv4InIPv6Prefix)
	copy(ip[12:16], sa.Addr[:])
	return ip
}
func _sockaddrInet6ToIPAndZone(sa *unix.SockaddrInet6) (net.IP, string) {
	ip := make([]byte, 16)
	copy(ip, sa.Addr[:])
	return ip, _ip6ZoneToString(int(sa.ZoneId))
}
func _ip6ZoneToString(zone int) string {
	if zone == 0 {
		return ""
	}
	if ifi, err := net.InterfaceByIndex(zone); err == nil {
		return ifi.Name
	}
	return _int2decimal(uint(zone))
}

// Convert int to decimal string.
func _int2decimal(i uint) string {
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

func _sockaddrToTCPOrUnixAddr(sa unix.Sockaddr) net.Addr {
	switch sa := sa.(type) {
	case *unix.SockaddrInet4:
		ip := _sockaddrInet4ToIP(sa)
		return &net.TCPAddr{IP: ip, Port: sa.Port}
	case *unix.SockaddrInet6:
		ip, zone := _sockaddrInet6ToIPAndZone(sa)
		return &net.TCPAddr{IP: ip, Port: sa.Port, Zone: zone}
	case *unix.SockaddrUnix:
		return &net.UnixAddr{Name: sa.Name, Net: "unix"}
	}
	return nil
}

func _tcpListen(addr string) (int, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return -1, GenErrNoERR_NUM("net.ResolveTCPAddr fail:", err)
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, unix.IPPROTO_TCP)
	if err != nil {
		return -1, GenErrNoERR_NUM("unix.Socket fail:", err)
	}

	sa4 := &unix.SockaddrInet4{Port: tcpAddr.Port}
	if tcpAddr.IP != nil {
		if len(tcpAddr.IP) == 16 {
			copy(sa4.Addr[:], tcpAddr.IP[12:16]) // copy last 4 bytes of slice to array
		} else {
			copy(sa4.Addr[:], tcpAddr.IP) // copy all bytes of slice to array
		}
	}
	err = unix.Bind(fd, sa4)
	if err != nil {
		return -1, GenErrNoERR_NUM("unix.Bind fail:", err)
	}

	err = unix.Listen(fd, 128)
	if err != nil {
		return -1, GenErrNoERR_NUM("unix.Listen fail:", err)
	}

	return fd, nil
}

func _tcpAccept(listenFD int) (connFD int, addr string, err error) {
	connFD, sa, err := unix.Accept(listenFD)
	if err != nil {
		if err == unix.EAGAIN {
			return -1, "", nil
		}
		return -1, "", GenErrNoERR_NUM("unix.Accept fail:", err)
	}
	remoteAddr := _sockaddrToTCPOrUnixAddr(sa)
	addr = remoteAddr.String()
	return
}

func _tcpConnectNoBlock(addr string)(fd int, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return -1, GenErrNoERR_NUM("net.ResolveTCPAddr fail:", err)
	}

	fd, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, unix.IPPROTO_TCP)
	if err != nil {
		return -1, GenErrNoERR_NUM("unix.Socket fail:", err)
	}

	sa4 := &unix.SockaddrInet4{Port: tcpAddr.Port}
	if tcpAddr.IP != nil {
		if len(tcpAddr.IP) == 16 {
			copy(sa4.Addr[:], tcpAddr.IP[12:16]) // copy last 4 bytes of slice to array
		} else {
			copy(sa4.Addr[:], tcpAddr.IP) // copy all bytes of slice to array
		}
	}
	err = unix.Connect(fd, sa4)
	if err != nil {
		if err != unix.EINPROGRESS {
			return -1, GenErrNoERR_NUM("connect fail:", err)
		}
	}

	return fd, nil
}

/*func TMP_tcpWrite(fd FD_DEF, bin []byte) (write_sz int, err error, bAgain bool) {
	return _tcpWrite(fd.FD, bin)
}*/
func _tcpWrite(fd int, bin []byte) (write_sz int, err error, bAgain bool) {
	n, err := unix.Write(fd, bin)
	if err != nil {
		if err != unix.EAGAIN {
			return 0, GenErrNoERR_NUM(" unix.Write fail, fd:", fd, " err:", err), false
		} else {
			bAgain = true
		}
	} else {
		bAgain = false
	}
	return n, nil, bAgain
}


func _tcpRead(fd int, bin []byte) (int, error) {
	n, err := unix.Read(fd, bin)
	if n == 0 {
		return 0, unix.ECONNRESET
	}
	if err != nil{
		if err == unix.EAGAIN {
			return 0, nil
		} else {
			return 0, err
		}
	}

	return n, err

	/*
		//from c++
		int32 ret = recv(fd, (int8 *)buf, (int32)buf_sz, 0);
		if(ret > 0)
			return ret;
		else if(0 == ret)
			return -1;
		else
		{
			if(EAGAIN == errno)
				return 0;
			return -1;
		}
	*/
}

func TcpGetPeerName(fd int) net.Addr {
	return _tcpGetPeerName(fd)
}
func _tcpGetPeerName(fd int) net.Addr {
	sa, err := unix.Getpeername(fd)
	if err != nil {
		return nil
	}
	return _sockaddrToTCPOrUnixAddr(sa)
}

func _setNoBlock(fd int) error{
	return unix.SetNonblock(fd, true)
}

func _setNoDelay(fd int) error {
	return unix.SetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_NODELAY, 1)
}

func _setDelay(fd int) error {
	return unix.SetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_NODELAY, 0)
}

func _setRecvBuffer(fd, size int) error {
	return unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_RCVBUF, size)
}

func _setSendBuffer(fd, size int) error {
	return unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_SNDBUF, size)
}

func _setReuseport(fd, reusePort int) error {
	return unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, reusePort)
}

func _setReuseAddr(fd, reuseAddr int) error {
	return unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, reuseAddr)
}

func _setLingerOff(fd int) error {
	l := &unix.Linger{Onoff:0,Linger:0}
	return unix.SetsockoptLinger(fd, syscall.SOL_SOCKET, syscall.SO_LINGER, l)
}
func _setLinger(fd, sec int) error {
	var l unix.Linger
	if sec >= 0 {
		l.Onoff = 1
		l.Linger = int32(sec)
	} else {
		l.Onoff = 0
		l.Linger = 0
	}
	return unix.SetsockoptLinger(fd, syscall.SOL_SOCKET, syscall.SO_LINGER, &l)
}

func _tcpKeepAlive(fd int, idle int, interval int, retry_count int) error {
	err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_KEEPALIVE, 1)
	if err != nil {
		return err
	}

	err = unix.SetsockoptInt(fd, unix.SOL_TCP, unix.TCP_KEEPIDLE, idle)
	if err != nil {
		unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_KEEPALIVE, 0)
		return err
	}

	err = unix.SetsockoptInt(fd, unix.SOL_TCP, unix.TCP_KEEPINTVL, interval)
	if err != nil {
		unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_KEEPALIVE, 0)
		return err
	}

	err = unix.SetsockoptInt(fd, unix.SOL_TCP, unix.TCP_KEEPCNT, retry_count)
	if err != nil {
		unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_KEEPALIVE, 0)
		return err
	}

	return nil
}
