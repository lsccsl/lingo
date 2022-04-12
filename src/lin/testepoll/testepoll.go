//go:build linux
// +build linux

package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"time"
)

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

func main() {
	efd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	fmt.Println(efd, " err:", err)

	fd := tcpListen("192.168.2.129:3001")

	conn_fd, sa, err := unix.Accept(fd)
	fmt.Println("new connect:", conn_fd, sa, err)

	for {
		time.Sleep(time.Second * 10)
	}
}
