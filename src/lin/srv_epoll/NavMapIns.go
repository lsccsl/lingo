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
	"sync"
	"unsafe"
)

type Coord3f struct {
	X float32
	Y float32
	Z float32
}

type nav_obstacle struct {
	center Coord3f
	half_ext Coord3f
	y_radian float32
}

type MAP_OBSTACLE map[uint32]*nav_obstacle

type NavMapIns struct {
	nav_lock_ sync.Mutex

	map_obstacle_ MAP_OBSTACLE

	handle_map_ins_ unsafe.Pointer
}

func ConstructNavMapIns() *NavMapIns {
	ins := &NavMapIns{
		map_obstacle_ : make(MAP_OBSTACLE),
	}

	return ins
}

func (pthis *NavMapIns)load_from_template(navMap *NavMap) {
	if pthis.handle_map_ins_ != nil {
		C.nav_delete(pthis.handle_map_ins_)
	}
	pthis.map_obstacle_ = make(MAP_OBSTACLE)

	pthis.handle_map_ins_ = unsafe.Pointer(C.nav_new())
	C.nav_load_from_template(pthis.handle_map_ins_, navMap.handle_template_)
}

func (pthis*NavMapIns)path_find(src * Coord3f, dst * Coord3f) (path []Coord3f){
	pthis.nav_lock_.Lock()
	defer pthis.nav_lock_.Unlock()

	var start_pos C.struct_RecastVec3f
	start_pos.x = C.float(src.X)
	start_pos.y = C.float(src.Y)
	start_pos.z = C.float(src.Z)
	var end_pos C.RecastPosT
	end_pos.x = C.float(dst.X)
	end_pos.y = C.float(dst.Y)
	end_pos.z = C.float(dst.Z)
	var pos *C.RecastPosT
	var pos_sz C.int
	C.nav_findpath(pthis.handle_map_ins_, &start_pos, &end_pos, &pos, &pos_sz, true)
	for i:=0; i < int(pos_sz); i ++ {
		tmp_v := uintptr(unsafe.Pointer(pos))  + uintptr(i)*unsafe.Sizeof(*pos)
		tmp_pos_ptr := (*C.RecastPosT)( unsafe.Pointer(tmp_v) )
		path = append(path, Coord3f{float32(tmp_pos_ptr.x), float32(tmp_pos_ptr.y), float32(tmp_pos_ptr.z)})
	}
	C.nav_freepath(pos)
	return
}



func (pthis*NavMapIns)add_obstacle(center * Coord3f, halfExtents * Coord3f, yRadians float32) (obstacle_id uint32) {
	pthis.nav_lock_.Lock()
	defer pthis.nav_lock_.Unlock()

	obstacle_id = uint32(C.nav_add_obstacle(pthis.handle_map_ins_, &C.struct_RecastVec3f{C.float(center.X), C.float(center.Y), C.float(center.Z)},
		&C.RecastPosT{C.float(halfExtents.X), C.float(halfExtents.Y), C.float(halfExtents.Z)},
		C.float(yRadians)))
	if obstacle_id == 0 {
		return
	}

	ob := &nav_obstacle{}
	ob.center = *center
	ob.half_ext = *halfExtents
	ob.y_radian = yRadians
	pthis.map_obstacle_[obstacle_id] = ob

	C.nav_update(pthis.handle_map_ins_)
	return
}

func (pthis*NavMapIns)del_obstacle(obstacle_id uint32)  {
	pthis.nav_lock_.Lock()
	defer pthis.nav_lock_.Unlock()

	delete(pthis.map_obstacle_, obstacle_id)

	C.nav_del_obstacle(pthis.handle_map_ins_, C.uint(obstacle_id))
	C.nav_update(pthis.handle_map_ins_)
}

func (pthis*NavMapIns)get_all_obstacle() MAP_OBSTACLE {
	pthis.nav_lock_.Lock()
	defer pthis.nav_lock_.Unlock()

	m := make(MAP_OBSTACLE)
	for k,v := range pthis.map_obstacle_ {
		m[k] = v
	}
	return m
}
