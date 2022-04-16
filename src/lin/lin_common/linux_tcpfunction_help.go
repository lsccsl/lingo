//go:build linux
// +build linux

package lin_common

import (
	"golang.org/x/sys/unix"
	"net"
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
		return -1, "", err
	}
	remoteAddr := _sockaddrToTCPOrUnixAddr(sa)
	addr = remoteAddr.String()
	return
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