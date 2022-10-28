package main

/*
#cgo CFLAGS: -I../../cpp/navwrapper
#cgo LDFLAGS: -L../../cpp/navwrapper/bin -lnavwrapper
#include "RecastCWrapper.h"
typedef void * VoidPtr;
typedef struct RecastVec3f RecastPosT;
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func test_recast_template() {

	str_file_path := "./test_mesh/test_scene.obj"

	bin_file_path := []byte(str_file_path)
	tmp := unsafe.Pointer(C.nav_temlate_create( (*C.char)( unsafe.Pointer(&bin_file_path[0]) ),
		6.0, 4.0, 4.0, 45.0))
	ins := unsafe.Pointer(C.nav_new())

	C.nav_load_from_template(ins, tmp)

	var start_pos C.struct_RecastVec3f
	start_pos.x = 702.190918
	start_pos.y = 1.53082275
	start_pos.z = 635.378662
	var end_pos C.RecastPosT
	end_pos.x = 710.805664
	end_pos.y = 1.00000000
	end_pos.z = 851.753296
	var pos *C.RecastPosT
	var pos_sz C.int
	C.nav_findpath(ins, &start_pos, &end_pos, &pos, &pos_sz, true)
	fmt.Println("pos_sz:", pos_sz)
	for i:=0; i < int(pos_sz); i ++ {
		fmt.Println(uintptr(i)*unsafe.Sizeof(*pos))
		tmp_v := uintptr(unsafe.Pointer(pos))  + uintptr(i)*unsafe.Sizeof(*pos)
		tmp_pos_ptr := (*C.RecastPosT)( unsafe.Pointer(tmp_v) )
		fmt.Println(tmp_pos_ptr.x, tmp_pos_ptr.y, tmp_pos_ptr.z)
	}
}

func main() {
	test_recast_template()
	fmt.Println("test_nav")

	str_file_path := "./test_mesh/nav_test.obj"

	bin_file_path := []byte(str_file_path)

	//ins1 := C.nav_create(C.CString("./test_mesh/nav_test.obj"))
	//fmt.Println("ins1 addr:", ins1)

	ins := unsafe.Pointer(C.nav_new())
	C.nav_reset_agent(ins, 2.0, 0.6, 0.9, 45.0)
	C.nav_load(ins, (*C.char)( unsafe.Pointer(&bin_file_path[0]) ) )
	fmt.Println("ins addr:", ins)

	id := C.nav_add_obstacle(ins, &C.struct_RecastVec3f{48.2378387, -1.40648651, 8.61733246},
		&C.RecastPosT{2.0, 2.0, 2.0},
		(45.0 / 360.0) * 2.0 * 3.14)
	C.nav_update(ins)
	C.nav_update(ins)
	fmt.Println("obstacle id:", id)

	var start_pos C.struct_RecastVec3f
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
		fmt.Println(uintptr(i)*unsafe.Sizeof(*pos))
		tmp_v := uintptr(unsafe.Pointer(pos))  + uintptr(i)*unsafe.Sizeof(*pos)
		tmp_pos_ptr := (*C.RecastPosT)( unsafe.Pointer(tmp_v) )
		fmt.Println(tmp_pos_ptr.x, tmp_pos_ptr.y, tmp_pos_ptr.z)
	}
	C.nav_freepath(pos)
}

/*
char -->  C.char -->  byte
signed char -->  C.schar -->  int8
unsigned char -->  C.uchar -->  uint8
short int -->  C.short -->  int16
short unsigned int -->  C.ushort -->  uint16
int -->  C.int -->  int
unsigned int -->  C.uint -->  uint32
long int -->  C.long -->  int32 or int64
long unsigned int -->  C.ulong -->  uint32 or uint64
long long int -->  C.longlong -->  int64
long long unsigned int -->  C.ulonglong -->  uint64
float -->  C.float -->  float32
double -->  C.double -->  float64
wchar_t -->  C.wchar_t  -->
void * -> unsafe.Pointer
*/