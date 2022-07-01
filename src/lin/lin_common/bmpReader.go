package lin_common

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

	BmpData []uint8
	BinPalette []uint8
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

func (pthis*Bitmap)WriteBmp(bmpFile string, binData[]uint8, w int, h int, bitCount int)error{
	file, err := os.Create(bmpFile)
	if err != nil {
		return err
	}

	defer func() {
		file.Sync()
		file.Close()
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
	err = binary.Write(file, binary.LittleEndian, uint32(headerLen + len(binData)))
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
	err = binary.Write(file, binary.LittleEndian, binData)
	if err != nil {
		return err
	}
	return nil
}

func (pthis*Bitmap)GetWidth() int {
	return int(pthis.Biheader.Width)
}

func (pthis*Bitmap)GetHeight() int {
	return int(pthis.Biheader.Height)
}