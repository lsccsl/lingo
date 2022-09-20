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
	"lin/lin_common"
	"unsafe"
)

type Coord3f struct {
	X float32
	Y float32
	Z float32
}

type NavMap struct {
	handle_nav_map_ unsafe.Pointer
}

func ConstructorNavMapMgr(path string) *NavMap {
	nav_map := &NavMap{}

	nav_map.handle_nav_map_ = C.nav_create(C.CString(path))

	return nav_map
}

func (pthis*NavMap)path_find(src Coord3f, dst Coord3f) (path []Coord3f){
	var start_pos C.struct_RecastPos
	start_pos.x = C.float(src.X)
	start_pos.y = C.float(src.Y)
	start_pos.z = C.float(src.Z)
	var end_pos C.RecastPosT
	end_pos.x = C.float(dst.X)
	end_pos.y = C.float(dst.Y)
	end_pos.z = C.float(dst.Z)
	var pos *C.RecastPosT
	var pos_sz C.int
	C.nav_findpath(pthis.handle_nav_map_, &start_pos, &end_pos, &pos, &pos_sz, true)
	lin_common.LogDebug("pos_sz:", pos_sz)
	for i:=0; i < int(pos_sz); i ++ {
		tmp_v := uintptr(unsafe.Pointer(pos))  + uintptr(i)*unsafe.Sizeof(*pos)
		tmp_pos_ptr := (*C.RecastPosT)( unsafe.Pointer(tmp_v) )
		lin_common.LogDebug(tmp_pos_ptr.x, tmp_pos_ptr.y, tmp_pos_ptr.z)
		path = append(path, Coord3f{float32(tmp_pos_ptr.x), float32(tmp_pos_ptr.y), float32(tmp_pos_ptr.z)})
	}
	C.nav_freepath(pos)
	return
}