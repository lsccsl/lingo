package main

import "lin/lin_common"

type MapAoiInf interface {
	Ntf_node_in_view(nodeID int, nodeIDInView int)
	Ntf_node_out_view(nodeID int, nodeIDOutView int)
}

type MapAOI struct {
	// cross-link
	crossLink *lin_common.Crosslink_mgr

	mapAoi map[int]MapAoiInf
}

func ConstructorMapAOI() *MapAOI {
	aoi := &MapAOI{mapAoi:make(map[int]MapAoiInf)}
	aoi.crossLink = lin_common.Crosslink_mgr_constructor(aoi)

	return aoi
}

func (pthis*MapAOI)Ntf_node_in_view(nodeID int, nodeIDInView int) {
	v, ok := pthis.mapAoi[nodeID]
	if !ok {
		return
	}
	v.Ntf_node_in_view(nodeID, nodeIDInView)
}
func (pthis*MapAOI)Ntf_node_out_view(nodeID int, nodeIDOutView int) {
	v, ok := pthis.mapAoi[nodeID]
	if !ok {
		return
	}
	v.Ntf_node_out_view(nodeID, nodeIDOutView)
}

func (pthis*MapAOI)add(X float32, Y float32, ViewRange float32, ntf MapAoiInf) int {
	aoiID := pthis.crossLink.Crosslink_mgr_add(&lin_common.Crosslink_node_param{X, Y, ntf, ViewRange})

	pthis.mapAoi[aoiID] = ntf

	return aoiID
}

func (pthis*MapAOI)del(aoiID int) {
	_, ok := pthis.mapAoi[aoiID]
	if !ok {
		return
	}

	pthis.crossLink.Crosslink_mgr_del(aoiID)
	delete(pthis.mapAoi, aoiID)
}

func (pthis*MapAOI)update(aoiID int, X float32, Y float32) {
	_, ok := pthis.mapAoi[aoiID]
	if !ok {
		return
	}

	pthis.crossLink.Crosslink_mgr_update_pos(aoiID, X, Y)
}
