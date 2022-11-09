package main

import "lin/lin_common"

type MapAOINtfClient struct {
	srvMgr *ServerMgr
	clientID int64
	objID int
	fd lin_common.FD_DEF

	mapView map[int]int
	mapViewBy map[int]int
}

func (pthis *MapAOINtfClient)Ntf_in_view(objID int) {
	lin_common.LogDebug(pthis.objID, " recv in view aoiID:", objID)

	pthis.mapView[objID] = objID
	//todo send msg to client


	//pthis.srvMgr.SendProtoMsg()
}

func (pthis *MapAOINtfClient)Ntf_out_view(objID int) {
	lin_common.LogDebug(pthis.objID, " recv out view aoiID:", objID)

	delete(pthis.mapView, objID)

	//todo send msg to client
}

func (pthis *MapAOINtfClient)Ntf_in_viewby(objID int) {
	lin_common.LogDebug(pthis.objID, " recv in viewby aoiID:", objID)

	pthis.mapViewBy[objID] = objID

	//todo send msg to client
}

func (pthis *MapAOINtfClient)Ntf_out_viewby(aoiID int) {
	lin_common.LogDebug(pthis.objID, " recv out viewby aoiID:", aoiID)

	delete(pthis.mapViewBy, aoiID)

	//todo send msg to client
}

func (pthis *MapAOINtfClient)setObjID(objID int) {
	pthis.objID = objID
}

func ConstructMapAOINtfClient(clientID int64, fd lin_common.FD_DEF, srvMgr *ServerMgr) *MapAOINtfClient {
	aoi := &MapAOINtfClient{
		clientID:clientID,
		srvMgr:srvMgr,
		fd:fd,
	}

	return aoi
}