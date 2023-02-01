package common

import (
	"encoding/binary"
	"os"
)

/*
typedef struct tagBITMAPFILEHEADER
{
UINT16 bfType;
DWORD bfSize;
UINT16 bfReserved1;
UINT16 bfReserved2;
DWORD bfOffBits;
} BITMAPFILEHEADER;
*/

const (
	BITMAP_BIN_HEADER int = 40
	BITMAP_HEADER_TOTAL int = 54
)

type BitmapFileHeader struct {
	BfType uint16 // BM
	BfSize uint32
	BfReserved1 uint16 // 0
	BfReserved2 uint16 // 0
	BfOffBits uint32
}//14bit

type BitmapInfoHeader struct {
	Size   uint32
	Width   int32
	Height   int32
	Places   uint16
	BitCount  uint16
	Compression uint32
	SizeImage  uint32
	XperlsPerMeter int32
	YperlsPerMeter int32
	ClsrUsed  uint32
	ClrImportant uint32
}//40bit

type Bitmap struct{
	BfHeader BitmapFileHeader
	Biheader BitmapInfoHeader

	widBytePitch int

	BmpData []uint8
	BinPalette []uint8
}

func CalWidBytePitch(wid int, bitCount int) int {
	widByte := (int(wid) * int(bitCount) + 7)/8
	return ((widByte + 3)/4) * 4
}

func (pthis*Bitmap)calWidBytePitch() {
	pthis.widBytePitch = CalWidBytePitch(int(pthis.Biheader.Width), int(pthis.Biheader.BitCount))
/*	widByte := (int(pthis.Biheader.Width) * int(pthis.Biheader.BitCount) + 7)/8
	pthis.widBytePitch = ((widByte + 3)/4) * 4*/
}

func (pthis*Bitmap)ReadBmp(bmpFile string)error{
	file, err := os.Open(bmpFile)
	if err != nil {
		return err
	}

	defer file.Close()

	//type拆成两个byte来读
	//Read第二个参数字节序一般windows/linux大部分都是LittleEndian,苹果系统用BigEndian
	//binary.Read(file, binary.LittleEndian, &pthis.headA)
	//binary.Read(file, binary.LittleEndian, &pthis.headB)
	err = binary.Read(file, binary.LittleEndian, &pthis.BfHeader.BfType)
	if err != nil {
		return err
	}

	//文件大小
	err = binary.Read(file, binary.LittleEndian, &pthis.BfHeader.BfSize)
	if err != nil {
		return err
	}

	//预留字节
	err = binary.Read(file, binary.LittleEndian, &pthis.BfHeader.BfReserved1)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &pthis.BfHeader.BfReserved2)
	if err != nil {
		return err
	}

	//偏移字节
	err = binary.Read(file, binary.LittleEndian, &pthis.BfHeader.BfOffBits)
	if err != nil {
		return err
	}

	err = binary.Read(file, binary.LittleEndian, &pthis.Biheader)
	if err != nil {
		return err
	}
	pthis.calWidBytePitch()

	szPalette := int(pthis.BfHeader.BfOffBits) - BITMAP_HEADER_TOTAL
	if szPalette > 0 {
		pthis.BinPalette = make([]uint8, szPalette)
		err = binary.Read(file, binary.LittleEndian, pthis.BinPalette)
		if err != nil {
			return err
		}
	}

	pthis.BmpData = make([]uint8, pthis.Biheader.SizeImage)
	err = binary.Read(file, binary.LittleEndian, pthis.BmpData)
	if err != nil {
		return err
	}
	return nil
}

func (pthis*Bitmap)WriteBmp(bmpFile string)error{
	file, err := os.Create(bmpFile)
	if err != nil {
		return err
	}

	defer func() {
		if file != nil {
			file.Sync()
			file.Close()
		}
	}()

	// write BM
	err = binary.Write(file, binary.LittleEndian, uint8('B'))
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.LittleEndian, uint8('M'))
	if err != nil {
		return err
	}

	headerLen := BITMAP_HEADER_TOTAL
	if pthis.BinPalette != nil {
		headerLen += len(pthis.BinPalette)
	}
	//write total file size
	err = binary.Write(file, binary.LittleEndian, uint32(headerLen + len(pthis.BmpData)))
	if err != nil {
		return err
	}

	//write reserve
	err = binary.Write(file, binary.LittleEndian, uint16(0))
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.LittleEndian, uint16(0))
	if err != nil {
		return err
	}

	//write offset
	err = binary.Write(file, binary.LittleEndian, uint32(headerLen))
	if err != nil {
		return err
	}

	//write bitmap info
	err = binary.Write(file, binary.LittleEndian, &pthis.Biheader)
	if err != nil {
		return err
	}

	//write palette
	err = binary.Write(file, binary.LittleEndian, pthis.BinPalette)
	if err != nil {
		return err
	}

	//write bitmap bin data
	err = binary.Write(file, binary.LittleEndian, pthis.BmpData)
	if err != nil {
		return err
	}
	return nil
}

func (pthis*Bitmap)GetRealWidth() int {
	return int(pthis.Biheader.Width)
}

func (pthis*Bitmap)GetPitchByteWidth() int {
	return pthis.widBytePitch
}

func (pthis*Bitmap)GetHeight() int {
	return int(pthis.Biheader.Height)
}

func CreateBMP(w int, h int, bitCount int, binBMP []uint8) *Bitmap {
	bmp := &Bitmap{}

	bmp.BfHeader.BfSize = uint32(BITMAP_HEADER_TOTAL + len(binBMP))
	bmp.BfHeader.BfOffBits = uint32(BITMAP_HEADER_TOTAL)

	bmp.Biheader.Size = uint32(BITMAP_BIN_HEADER)
	bmp.Biheader.Width = int32(w)
	bmp.Biheader.Height = int32(h)
	bmp.Biheader.Places = 1
	bmp.Biheader.BitCount = uint16(bitCount)
	bmp.Biheader.Compression = 0
	bmp.Biheader.SizeImage = uint32(len(binBMP))
	bmp.Biheader.XperlsPerMeter = 0
	bmp.Biheader.YperlsPerMeter = 0
	bmp.Biheader.ClsrUsed = 0
	bmp.Biheader.ClrImportant = 0

	bmp.calWidBytePitch()

	if binBMP != nil {
		bmp.BmpData = make([]uint8, len(binBMP))
		copy(bmp.BmpData, binBMP)
	}

	return bmp
}
