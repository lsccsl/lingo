package main

import (
	"lin/lin_common"
	"lin/navmeshwrapper"
	"math"
	"math/rand"
	"time"
)

type MapAoiTestView struct {
	srvMgr *ServerMgr
	objID int

	mapView map[int]int
	mapViewBy map[int]int
}

func (pthis*MapAoiTestView)Ntf_in_view(objID int) {
	lin_common.LogDebug(pthis.objID, " aoiID:", objID)

	pthis.mapView[objID] = objID
}
func (pthis*MapAoiTestView)Ntf_out_view(objID int) {
	lin_common.LogDebug(pthis.objID, " aoiID:", objID)

	delete(pthis.mapView, objID)
}
func (pthis*MapAoiTestView)Ntf_in_viewby(objID int) {
	lin_common.LogDebug(pthis.objID, " aoiID:", objID)

	pthis.mapViewBy[objID] = objID
}
func (pthis*MapAoiTestView)Ntf_out_viewby(objID int) {
	lin_common.LogDebug(pthis.objID, " aoiID:", objID)

	delete(pthis.mapViewBy, objID)
}
func (pthis*MapAoiTestView)setObjID(objID int) {
	pthis.objID = objID
}

func (pthis*MapProcess)initTestViewNode() {
	rSeed := time.Now().Unix()
	lin_common.LogDebug("rSeed:", rSeed)
	rand.Seed(rSeed)

	boundMin := &navmeshwrapper.Coord3f{}
	boundMax := &navmeshwrapper.Coord3f{}
	pthis.navIns.GetNavBound(boundMin, boundMax)

	boundMin.X = float32(math.Ceil(float64(boundMin.X)))
	boundMin.Y = float32(math.Ceil(float64(boundMin.Y)))
	boundMin.Z = float32(math.Ceil(float64(boundMin.Z)))
	boundMax.X = float32(math.Ceil(float64(boundMax.X)))
	boundMax.Y = float32(math.Ceil(float64(boundMax.Y)))
	boundMax.Z = float32(math.Ceil(float64(boundMax.Z)))

	xRange := int(boundMax.X - boundMin.X)
	zRange := int(boundMax.Z - boundMin.Z)
	lin_common.LogDebug("xRange:", xRange, " zRange:", zRange)

	for i := 0; i < 100; i ++ {
		aoi := &MapAoiTestView{
			srvMgr : pthis.procMgr.eSrvMgr,
			mapView : make(map[int]int),
			mapViewBy : make(map[int]int),
		}
		x := boundMin.X + float32(rand.Int() % xRange)
		z := boundMin.Z + float32(rand.Int() % zRange)
		aoi.objID = pthis.aoi.genID()
		pthis.aoi.add(aoi.objID, x, z, 50, aoi)
		lin_common.LogDebug("add node:", aoi.objID, "[", x, "-", z, "]")
	}

	{
		aoi := &MapAoiTestView{
			srvMgr:    pthis.procMgr.eSrvMgr,
			mapView:   make(map[int]int),
			mapViewBy: make(map[int]int),
		}
		objID := pthis.aoi.genID()
		lin_common.LogDebug("add node:", objID)
		pthis.aoi.add(objID, 0, 0, 50, aoi)
	}
}