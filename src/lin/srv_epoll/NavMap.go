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
	"lin/lin_common"
	"sync"
	"unsafe"
)


type NavMap struct {
	nav_lock_ sync.Mutex
	handle_template_ unsafe.Pointer
}

func ConstructorNavMapMgr(file_path string) *NavMap {

	nav_map := &NavMap{}

	{
		nav_map.handle_template_ = unsafe.Pointer(C.nav_temlate_create(C.CString(file_path),
			6.0, 4.0, 4.0, 45.0))
		lin_common.LogDebug("load template success", file_path)

		navIns := ConstructNavMapIns()
		navIns.load_from_template(nav_map)
		lin_common.LogDebug("load map from template success", file_path)

		src := Coord3f{123.61628, 0, 101.47595}
		dst := Coord3f{966.7898,  0, 730.6272}
		path := navIns.path_find(&src, &dst)
		lin_common.LogDebug(len(path), " nav instance path:", path)
	}

	return nav_map
}
