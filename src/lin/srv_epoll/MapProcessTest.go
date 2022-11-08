package main

import (
	"lin/lin_common"
	"math"
	"math/rand"
	"time"
)

type MapAoiTestView struct {
	aoiID int
}

func (pthis*MapAoiTestView)Ntf_in_view(aoiID int) {
	lin_common.LogDebug(pthis.aoiID, " aoiID:", aoiID)
}
func (pthis*MapAoiTestView)Ntf_out_view(aoiID int) {
	lin_common.LogDebug(pthis.aoiID, " aoiID:", aoiID)
}
func (pthis*MapAoiTestView)Ntf_in_viewby(aoiID int) {
	lin_common.LogDebug(pthis.aoiID, " aoiID:", aoiID)
}
func (pthis*MapAoiTestView)Ntf_out_viewby(aoiID int) {
	lin_common.LogDebug(pthis.aoiID, " aoiID:", aoiID)
}
func (pthis*MapAoiTestView)setAOIID(aoiID int) {
	pthis.aoiID = aoiID
}

func (pthis*MapProcess)initTestViewNode() {
	rSeed := time.Now().Unix()
	lin_common.LogDebug("rSeed:", rSeed)
	rand.Seed(rSeed)

	boundMin := &Coord3f{}
	boundMax := &Coord3f{}
	pthis.navIns.getNavBound(boundMin, boundMax)

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
		aoi := &MapAoiTestView{}
		x := boundMin.X + float32(rand.Int() % xRange)
		z := boundMin.Z + float32(rand.Int() % zRange)
		aoi.aoiID = pthis.aoi.add(x, z, 10, aoi)
		lin_common.LogDebug("add node:", aoi.aoiID, "[", x, "-", z, "]")
	}
}