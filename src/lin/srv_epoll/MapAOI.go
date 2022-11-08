package main

import "lin/lin_common"

type MapAoiInf interface {
	Ntf_in_view(aoiID int)
	Ntf_out_view(aoiID int)
	Ntf_in_viewby(aoiID int)
	Ntf_out_viewby(aoiID int)
	setAOIID(aoiID int)
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
	node, _ := pthis.mapAoi[nodeID]
	if node != nil {
		node.Ntf_in_view(nodeIDInView)
	}
	nodeView, _ := pthis.mapAoi[nodeIDInView]
	if nodeView != nil {
		nodeView.Ntf_in_viewby(nodeID)
	}
}
func (pthis*MapAOI)Ntf_node_out_view(nodeID int, nodeIDOutView int) {
	node, _ := pthis.mapAoi[nodeID]
	if node != nil {
		node.Ntf_out_view(nodeIDOutView)
	}
	nodeView, _ := pthis.mapAoi[nodeIDOutView]
	if nodeView != nil {
		nodeView.Ntf_out_viewby(nodeID)
	}
}

func (pthis*MapAOI)add(X float32, Y float32, ViewRange float32, ntf MapAoiInf) int {
	aoiID := pthis.crossLink.Crosslink_mgr_gen_id()
	ntf.setAOIID(aoiID)

	pthis.mapAoi[aoiID] = ntf
	pthis.crossLink.Crosslink_mgr_add(&lin_common.Crosslink_node_param{aoiID, X, Y, ntf, ViewRange})
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
