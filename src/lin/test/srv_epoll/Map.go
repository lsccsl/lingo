package main

import (
	"lin/lin_common"
	"lin/lin_common/pathsearch"
	"lin/msgpacket"
)

type MapMgr struct {
	mapData *pathsearch.MapData
}

func ConstructorMapMgr(path string) *MapMgr {
	mapMgr := &MapMgr{}

	mapMgr.mapData = &pathsearch.MapData{}
	mapMgr.mapData.LoadMap(path)

	lin_common.LogDebug("wid:", mapMgr.mapData.GetWidReal(),
		" hei:", mapMgr.mapData.GetHeight(),
		" pitch:", mapMgr.mapData.GetWidPitch())

	return mapMgr
}

func (pthis*MapMgr)GetMapProtoMsg(msg *msgpacket.MSG_GET_MAP_RES) {
	msg.MapWid = int32(pthis.mapData.GetWidReal())
	msg.MapHigh = int32(pthis.mapData.GetHeight())
	msg.MapPitch = int32(pthis.mapData.GetWidPitch())
	msg.MapData = *pthis.mapData.GetMapBit()
}
