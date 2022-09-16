package main

/*
#cgo CFLAGS: -I../../cpp/navwrapper
#cgo LDFLAGS: -L../../cpp/navwrapper/bin -lnavwrapper
#include "RecastCWrapper.h"
typedef void * VoidPtr;
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("test_nav")

	str_file_path := "./test_mesh/nav_test.obj"

	bin_file_path := []byte(str_file_path)

	ins := unsafe.Pointer(C.nav_create((*C.char)( unsafe.Pointer(&bin_file_path[0]) ) ) )
	fmt.Println("ins addr:", ins)

	//C.nav_findpath(ins)
}
