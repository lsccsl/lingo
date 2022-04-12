//go:build linux
// +build linux

package main

import (
	"fmt"
	"golang.org/x/sys/unix"
)

func main() {
	efd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	fmt.Println(efd, " err:", err)
}