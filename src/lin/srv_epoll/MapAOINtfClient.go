package main

import "lin/lin_common"

type MapAOINtfClient struct {
	srvMgr *ServerMgr
	clientID int64
	aoiID int

	mapView map[int]int
	mapViewBy map[int]int
}

func (pthis *MapAOINtfClient)Ntf_in_view(aoiID int) {
	lin_common.LogDebug(pthis.aoiID, " recv in view aoiID:", aoiID)

	pthis.mapView[aoiID] = aoiID
	//todo send msg to client
}

func (pthis *MapAOINtfClient)Ntf_out_view(aoiID int) {
	lin_common.LogDebug(pthis.aoiID, " recv out view aoiID:", aoiID)

	delete(pthis.mapView, aoiID)

	//todo send msg to client
}

func (pthis *MapAOINtfClient)Ntf_in_viewby(aoiID int) {
	lin_common.LogDebug(pthis.aoiID, " recv in viewby aoiID:", aoiID)

	pthis.mapViewBy[aoiID] = aoiID

	//todo send msg to client
}

func (pthis *MapAOINtfClient)Ntf_out_viewby(aoiID int) {
	lin_common.LogDebug(pthis.aoiID, " recv out viewby aoiID:", aoiID)

	delete(pthis.mapViewBy, aoiID)

	//todo send msg to client
}

func (pthis *MapAOINtfClient)setAOIID(aoiID int) {
	pthis.aoiID = aoiID
}

func ConstructMapAOINtfClient(clientID int64, srvMgr *ServerMgr) *MapAOINtfClient {
	aoi := &MapAOINtfClient{
		clientID:clientID,
		srvMgr:srvMgr,
	}

	return aoi
}