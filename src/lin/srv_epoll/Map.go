package main

import (
	"lin/lin_common"
	"lin/msgpacket"
)

type MapMgr struct {
	mapData *lin_common.MapData
}

func ConstructorMapMgr(path string) *MapMgr {
	mapMgr := &MapMgr{}

	mapMgr.mapData = &lin_common.MapData{}
	mapMgr.mapData.LoadMap(path)

	return mapMgr
}

func (pthis*MapMgr)GetMapProtoMsg(msg *msgpacket.MSG_GET_MAP_RES) {
	msg.MapWid = int32(pthis.mapData.GetWidReal())
	msg.MapHigh = int32(pthis.mapData.GetHeight())
	msg.MapPitch = int32(pthis.mapData.GetWidPitch())
	msg.MapData = *pthis.mapData.GetMapBit()
}
