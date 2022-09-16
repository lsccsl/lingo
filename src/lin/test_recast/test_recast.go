package main

/*
#cgo CFLAGS: -I../../cpp/navwrapper
#cgo LDFLAGS: -L../../cpp/navwrapper/bin -lnavwrapper
#include "RecastCWrapper.h"
typedef void * VoidPtr;
typedef struct RecastPos RecastPosT;
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

	ins1 := C.nav_create(C.CString("./test_mesh/nav_test.obj"))
	fmt.Println("ins1 addr:", ins1)

	ins := unsafe.Pointer(C.nav_create((*C.char)( unsafe.Pointer(&bin_file_path[0]) ) ) )
	fmt.Println("ins addr:", ins)

	var start_pos C.struct_RecastPos
	start_pos.x = 40.5650635
	start_pos.y = -1.71816540
	start_pos.z = 22.0546188
	var end_pos C.RecastPosT
	end_pos.x = 49.6740074
	end_pos.y = -2.50520134
	end_pos.z = -6.56286621
	var pos *C.RecastPosT
	var pos_sz C.int
	C.nav_findpath(ins, &start_pos, &end_pos, &pos, &pos_sz, true)
	fmt.Println("pos_sz:", pos_sz)
	for i:=0; i < int(pos_sz); i ++ {
	}
}
