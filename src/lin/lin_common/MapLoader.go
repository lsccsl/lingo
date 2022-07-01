package lin_common

type MapData struct {
	wid int
	hei int

	mapBit []uint8
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

func (pthis*MapData)GetBitBlock(x int, y int)bool{
	if x < 0 || x >= pthis.wid {
		return true
	}
	if y < 0 || y >= pthis.hei {
		return true
	}

	idx := x * pthis.wid + y
	idxByte := idx / 8
	idxBit := 7 - idx % 8
	posByte := pthis.mapBit[idxByte]
	posBit := posByte & (1 << idxBit)

	return posBit != 0
}

func (pthis*MapData)DumpMap(mapFile string){
}