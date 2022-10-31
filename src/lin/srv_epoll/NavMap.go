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

type NavMap struct {
	nav_lock_ sync.Mutex

	handle_nav_map_ unsafe.Pointer
	map_obstacle_ MAP_OBSTACLE

	handle_map_ins_ unsafe.Pointer
	handle_template_ unsafe.Pointer
}

func ConstructorNavMapMgr(file_path string) *NavMap {

	nav_map := &NavMap{}

	nav_map.map_obstacle_ = make(MAP_OBSTACLE)

	nav_map.handle_nav_map_ = C.nav_create(C.CString(file_path))
	if nav_map.handle_nav_map_ == nil {
		lin_common.LogErr("fail load", file_path)
		return nil
	}
	lin_common.LogDebug("load success", file_path)

/*	{
		src := Coord3f{123.61628, 0, 101.47595}
		dst := Coord3f{966.7898,  0, 730.6272}
		path := nav_map.path_find(&src, &dst)
		lin_common.LogDebug(len(path), " path:", path)
	}*/

	{
		nav_map.handle_template_ = unsafe.Pointer(C.nav_temlate_create(C.CString(file_path),
			6.0, 4.0, 4.0, 45.0))
		nav_map.handle_map_ins_ = unsafe.Pointer(C.nav_new())
		C.nav_load_from_template(nav_map.handle_map_ins_, nav_map.handle_template_)

		src := Coord3f{123.61628, 0, 101.47595}
		dst := Coord3f{966.7898,  0, 730.6272}
		path := nav_map.path_find(&src, &dst)
		lin_common.LogDebug(len(path), " nav instance path:", path)
	}

	return nav_map
}

/*func (pthis*NavMap)path_find(src * Coord3f, dst * Coord3f) (path []Coord3f){
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
	C.nav_findpath(pthis.handle_nav_map_, &start_pos, &end_pos, &pos, &pos_sz, true)
	for i:=0; i < int(pos_sz); i ++ {
		tmp_v := uintptr(unsafe.Pointer(pos))  + uintptr(i)*unsafe.Sizeof(*pos)
		tmp_pos_ptr := (*C.RecastPosT)( unsafe.Pointer(tmp_v) )
		path = append(path, Coord3f{float32(tmp_pos_ptr.x), float32(tmp_pos_ptr.y), float32(tmp_pos_ptr.z)})
	}
	C.nav_freepath(pos)
	return
}*/

func (pthis*NavMap)path_find(src * Coord3f, dst * Coord3f) (path []Coord3f){
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



func (pthis*NavMap)add_obstacle(center * Coord3f, halfExtents * Coord3f, yRadians float32) (obstacle_id uint32) {
	pthis.nav_lock_.Lock()
	defer pthis.nav_lock_.Unlock()

	obstacle_id = uint32(C.nav_add_obstacle(pthis.handle_nav_map_, &C.struct_RecastVec3f{C.float(center.X), C.float(center.Y), C.float(center.Z)},
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

	C.nav_update(pthis.handle_nav_map_)
	return
}

func (pthis*NavMap)del_obstacle(obstacle_id uint32)  {
	pthis.nav_lock_.Lock()
	defer pthis.nav_lock_.Unlock()

	delete(pthis.map_obstacle_, obstacle_id)

	C.nav_del_obstacle(pthis.handle_nav_map_, C.uint(obstacle_id))
	C.nav_update(pthis.handle_nav_map_)
}

func (pthis*NavMap)get_all_obstacle() MAP_OBSTACLE {
	pthis.nav_lock_.Lock()
	defer pthis.nav_lock_.Unlock()

	return pthis.map_obstacle_
}
