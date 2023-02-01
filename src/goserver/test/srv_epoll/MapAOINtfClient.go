package main

import (
	"goserver/common"
	"goserver/msgpacket"
)

type MapAOINtfClient struct {
	srvMgr *ServerMgr
	clientID int64
	objID int
	fd    common.FD_DEF

	mapView map[int]int
	mapViewBy map[int]int
}

func (pthis *MapAOINtfClient)Ntf_in_view(objID int) {
	common.LogDebug(pthis.objID, " recv in view aoiID:", objID)

	pthis.mapView[objID] = objID

	msg := &msgpacket.MSG_NTF_IN_VIEW{}
	msg.ObjId = int64(objID)
	pthis.srvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_NTF_IN_VIEW, msg)
}

func (pthis *MapAOINtfClient)Ntf_out_view(objID int) {
	common.LogDebug(pthis.objID, " recv out view aoiID:", objID)

	delete(pthis.mapView, objID)

	msg := &msgpacket.MSG_NTF_OUT_VIEW{}
	msg.ObjId = int64(objID)
	pthis.srvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_NTF_IN_VIEW, msg)
}

func (pthis *MapAOINtfClient)Ntf_in_viewby(objID int) {
	common.LogDebug(pthis.objID, " recv in viewby aoiID:", objID)

	pthis.mapViewBy[objID] = objID
}

func (pthis *MapAOINtfClient)Ntf_out_viewby(aoiID int) {
	common.LogDebug(pthis.objID, " recv out viewby aoiID:", aoiID)

	delete(pthis.mapViewBy, aoiID)
}

func (pthis *MapAOINtfClient)setObjID(objID int) {
	pthis.objID = objID
}

func ConstructMapAOINtfClient(clientID int64, fd common.FD_DEF, srvMgr *ServerMgr) *MapAOINtfClient {
	aoi := &MapAOINtfClient{
		clientID:clientID,
		srvMgr:srvMgr,
		fd:fd,
	}

	return aoi
}