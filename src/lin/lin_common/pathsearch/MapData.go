package pathsearch

import (
	"lin/lin_common"
	"math"
)



type Coord2d struct {
	X int
	Y int
}

func (pthis*Coord2d)Add(r*Coord2d) Coord2d {
	return Coord2d{pthis.X + r.X, pthis.Y + r.Y}
}
func (pthis*Coord2d)Dec(r*Coord2d) Coord2d {
	return Coord2d{pthis.X - r.X, pthis.Y - r.Y}
}
func (pthis*Coord2d)IsEqual(r*Coord2d) bool {
	return pthis.X ==  r.X && pthis.Y == r.Y
}
func (pthis*Coord2d)IsNear(r*Coord2d) bool {
	return math.Abs(float64(pthis.X -  r.X)) <= 1 && math.Abs(float64(pthis.Y - r.Y)) <= 1
}

type MapData struct {
	widReal int
	widBytePitch int
	hei int

	mapBit []uint8
}

func (pthis*MapData)LoadMap(mapFile string)error{
	bmp := lin_common.Bitmap{}
	err := bmp.ReadBmp(mapFile)
	if err != nil {
		return err
	}

	pthis.mapBit = bmp.BmpData
	pthis.widReal = bmp.GetRealWidth()
	pthis.widBytePitch = bmp.GetPitchByteWidth()
	pthis.hei = bmp.GetHeight()

	return nil
}

func (pthis*MapData)CoordIsBlock(pos Coord2d)bool{
	return pthis.IsBlock(pos.X, pos.Y)
}
func (pthis*MapData)IsBlock(x int, y int)bool{
	if x < 0 || x >= pthis.widReal {
		return true
	}
	if y < 0 || y >= pthis.hei {
		return true
	}

	idxByte := y * pthis.widBytePitch + x/8
	idxBit := 7 - x % 8
	posByte := pthis.mapBit[idxByte]
	posBit := posByte & (1 << idxBit)

	return posBit == 0
}

func (pthis*MapData)GetWidReal()int{
	return pthis.widReal
}
func (pthis*MapData)GetWidPitch()int{
	return pthis.widBytePitch
}
func (pthis*MapData)GetHeight()int{
	return pthis.hei
}
func (pthis*MapData)GetMapBit()*[]uint8{
	return &pthis.mapBit
}


func (pthis*MapData)DumpMap(strMapFile string, path []Coord2d, src *Coord2d, dst *Coord2d, searchMgr *SearchMgr) {

	widBytePitch := lin_common.CalWidBytePitch(pthis.widReal, 24)
	dataLen := widBytePitch * pthis.hei
	tmpBMP := make([]uint8, dataLen)

	for j := 0; j < pthis.hei; j ++ {
		for i := 0; i < pthis.widReal; i ++ {
			var clr uint8 = 0xff
			if pthis.IsBlock(i, j) {
				clr = 0
			}
			newIdx := j * widBytePitch + i * 3
			tmpBMP[newIdx + 0] = clr
			tmpBMP[newIdx + 1] = clr
			tmpBMP[newIdx + 2] = clr
		}
	}

	if searchMgr != nil {
		for key, _ := range searchMgr.mapHistoryPath {
			idx := key.Y*widBytePitch + key.X * 3
			tmpBMP[idx+0] = 0
			tmpBMP[idx+1] = 0xff
			tmpBMP[idx+2] = 0
		}
	}

	if path != nil {
		for _, val := range path {
			idx := val.Y*widBytePitch + val.X * 3
			tmpBMP[idx+0] = 0
			tmpBMP[idx+1] = 0
			tmpBMP[idx+2] = 0xff
		}
	}

	if src != nil {
		idx := src.Y*widBytePitch + src.X * 3
		tmpBMP[idx+0] = 0xff
		tmpBMP[idx+1] = 0
		tmpBMP[idx+2] = 0
	}
	if dst != nil {
		idx := dst.Y*widBytePitch + dst.X * 3
		tmpBMP[idx+0] = 0xff
		tmpBMP[idx+1] = 0
		tmpBMP[idx+2] = 0
	}

	bmp := lin_common.CreateBMP(pthis.widReal, pthis.hei, 24, tmpBMP)

	bmp.WriteBmp(strMapFile)
}
