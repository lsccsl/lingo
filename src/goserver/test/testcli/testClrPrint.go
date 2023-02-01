package main

import "syscall"

var kernelDll *syscall.LazyDLL = nil
var SetConsoleTextAttributePtr *syscall.LazyProc = nil
var CloseHandlePtr *syscall.LazyProc = nil

func BeginClrPrint(clr int){
	if nil == kernelDll {
		kernelDll = syscall.NewLazyDLL("kernel32.dll")
	}
	if nil == SetConsoleTextAttributePtr {
		SetConsoleTextAttributePtr = kernelDll.NewProc("SetConsoleTextAttribute")
	}
	handle, _, _ := SetConsoleTextAttributePtr.Call(uintptr(syscall.Stdout), uintptr(clr))
	if nil == CloseHandlePtr {
		CloseHandlePtr = kernelDll.NewProc("CloseHandle")
	}
	CloseHandlePtr.Call(handle)
}

func EndClrPrint(){
	if nil == kernelDll {
		kernelDll = syscall.NewLazyDLL("kernel32.dll")
	}
	if nil == SetConsoleTextAttributePtr {
		SetConsoleTextAttributePtr = kernelDll.NewProc("SetConsoleTextAttribute")
	}
	//kernel32 := syscall.NewLazyDLL("kernel32.dll")
	//proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := SetConsoleTextAttributePtr.Call(uintptr(syscall.Stdout), uintptr(7))
	if nil == CloseHandlePtr {
		CloseHandlePtr = kernelDll.NewProc("CloseHandle")
	}
	CloseHandlePtr.Call(handle)
}