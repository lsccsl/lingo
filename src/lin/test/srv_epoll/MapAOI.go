package main

import (
	"lin/lin_common/crosslink"
)

type MapAoiInf interface {
	Ntf_in_view(ObjID int)
	Ntf_out_view(ObjID int)
	Ntf_in_viewby(ObjID int)
	Ntf_out_viewby(ObjID int)

	setObjID(ObjID int)
}

type MapAOI struct {
	// cross-link
	crossLink *crosslink.Crosslink_mgr

	mapAoi map[int]MapAoiInf
}

func ConstructorMapAOI() *MapAOI {
	aoi := &MapAOI{mapAoi: make(map[int]MapAoiInf)}
	aoi.crossLink = crosslink.Crosslink_mgr_constructor(aoi)

	return aoi
}

func (pthis*MapAOI)Ntf_node_in_view(objID int, objIDInView int) {
	node, _ := pthis.mapAoi[objID]
	if node != nil {
		node.Ntf_in_view(objIDInView)
	}
	nodeView, _ := pthis.mapAoi[objIDInView]
	if nodeView != nil {
		nodeView.Ntf_in_viewby(objID)
	}
}
func (pthis*MapAOI)Ntf_node_out_view(objID int, objIDOutView int) {
	node, _ := pthis.mapAoi[objID]
	if node != nil {
		node.Ntf_out_view(objIDOutView)
	}
	nodeView, _ := pthis.mapAoi[objIDOutView]
	if nodeView != nil {
		nodeView.Ntf_out_viewby(objID)
	}
}

func (pthis*MapAOI)genID() int {
	return pthis.crossLink.Crosslink_mgr_gen_id()
}

func (pthis*MapAOI)add(objID int, X float32, Y float32, ViewRange float32, ntf MapAoiInf) {
	ntf.setObjID(objID)
	pthis.mapAoi[objID] = ntf
	pthis.crossLink.Crosslink_mgr_add(&crosslink.Crosslink_node_param{objID, X, Y, ntf, ViewRange})
}

func (pthis*MapAOI)del(objID int) {
	_, ok := pthis.mapAoi[objID]
	if !ok {
		return
	}

	pthis.crossLink.Crosslink_mgr_del(objID)
	delete(pthis.mapAoi, objID)
}

func (pthis*MapAOI)update(objID int, X float32, Y float32) {
	_, ok := pthis.mapAoi[objID]
	if !ok {
		return
	}

	pthis.crossLink.Crosslink_mgr_update_pos(objID, X, Y)
}
