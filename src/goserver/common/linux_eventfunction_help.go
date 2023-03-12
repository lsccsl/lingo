//go:build linux
// +build linux

package common

import (
	"golang.org/x/sys/unix"
)

func unixEvent() (int, error){
	evtFD, err := unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC)
	if err != nil {
		return -1, GenErrNoERR_NUM("unix.Eventfd err:", err)
	}
	return evtFD, nil
}