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

func (pthis*MapData)DumpMap(strMapFile string) {

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
			tmpBMP[newIdx + 1] = 0
			tmpBMP[newIdx + 2] = clr
		}
	}

	bmp := CreateBMP(pthis.wid, pthis.hei, 24, tmpBMP)

	bmp.WriteBmp(strMapFile)
}
