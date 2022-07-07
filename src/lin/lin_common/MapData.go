package lin_common


type Coord2d struct {
	X int
	Y int
}
type MapData struct {
	wid int
	hei int

	mapBit []uint8

	openNodeMgr SearchOpenNodeMgr
}

func (pthis*MapData)LoadMap(mapFile string)error{
	bmp := Bitmap{}
	err := bmp.ReadBmp(mapFile)
	if err != nil {
		return err
	}

	pthis.mapBit = bmp.BmpData
	pthis.wid = bmp.GetWidth()
	pthis.hei = bmp.GetHeight()

	return nil
}

func (pthis*MapData)CoordIsBlock(pos Coord2d)bool{
	return pthis.IsBlock(pos.X, pos.Y)
}
func (pthis*MapData)IsBlock(x int, y int)bool{
	if x < 0 || x >= pthis.wid {
		return true
	}
	if y < 0 || y >= pthis.hei {
		return true
	}

	idx := y * pthis.wid + x
	idxByte := idx / 8
	idxBit := 7 - idx % 8
	posByte := pthis.mapBit[idxByte]
	posBit := posByte & (1 << idxBit)

	return posBit == 0
}




func (pthis*MapData)DumpMap(strMapFile string, path []Coord2d, src * Coord2d , dst * Coord2d) {

	dataLen := len(pthis.mapBit)
	tmpBMP := make([]uint8, dataLen * 24)

	for idx, val := range pthis.mapBit {
		for i := 7; i >= 0 ; i -- {
			tmp := val & (1 << i)
			var clr uint8 = 0
			if tmp != 0 {
				clr = 0xff
			}
			newIdx := idx * 24 + (7 - i) * 3
			tmpBMP[newIdx + 0] = clr
			tmpBMP[newIdx + 1] = clr
			tmpBMP[newIdx + 2] = clr
		}
	}

	if path != nil {
		for _, val := range path {
			idx := (val.Y*pthis.wid + val.X) * 3
			tmpBMP[idx+0] = 0
			tmpBMP[idx+1] = 0
			tmpBMP[idx+2] = 0xff
		}
	}

	if src != nil {
		idx := (src.Y*pthis.wid + src.X) * 3
		tmpBMP[idx+0] = 0xff
		tmpBMP[idx+1] = 0
		tmpBMP[idx+2] = 0
	}
	if dst != nil {
		idx := (dst.Y*pthis.wid + dst.X) * 3
		tmpBMP[idx+0] = 0xff
		tmpBMP[idx+1] = 0
		tmpBMP[idx+2] = 0
	}

	bmp := CreateBMP(pthis.wid, pthis.hei, 24, tmpBMP)

	bmp.WriteBmp(strMapFile)
}
