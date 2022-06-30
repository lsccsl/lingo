package lin_common

import (
	"encoding/binary"
	"os"
	//"image"
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
	bfType uint16 // BM
	bfSize uint32
	bfReserved1 uint16 // 0
	bfReserved2 uint16 // 0
	bfOffBits uint32
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
	bfHeader BitmapFileHeader
	Biheader BitmapInfoHeader

	BmpData []uint8
	BinPalette []uint8
}

func (pthis*Bitmap)ReadBmp(bmpFile string)error{
	file, err := os.Open(bmpFile)
	if err != nil {
		LogDebug(err)
		return err
	}

	defer file.Close()

	//type拆成两个byte来读
	//Read第二个参数字节序一般windows/linux大部分都是LittleEndian,苹果系统用BigEndian
	//binary.Read(file, binary.LittleEndian, &pthis.headA)
	//binary.Read(file, binary.LittleEndian, &pthis.headB)
	err = binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfType)
	if err != nil {
		LogDebug(err)
		return err
	}

	//文件大小
	err = binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfSize)
	if err != nil {
		LogDebug(err)
		return err
	}

	//预留字节
	err = binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfReserved1)
	if err != nil {
		LogDebug(err)
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfReserved2)
	if err != nil {
		LogDebug(err)
		return err
	}

	//偏移字节
	err = binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfOffBits)
	if err != nil {
		LogDebug(err)
		return err
	}

	err = binary.Read(file, binary.LittleEndian, &pthis.Biheader)
	if err != nil {
		LogDebug(err)
		return err
	}

	szPalette := int(pthis.bfHeader.bfOffBits) - BITMAP_HEADER_TOTAL
	if szPalette > 0 {
		pthis.BinPalette = make([]uint8, szPalette)
		err = binary.Read(file, binary.LittleEndian, pthis.BinPalette)
		if err != nil {
			LogDebug(err)
			return err
		}
	}

	pthis.BmpData = make([]uint8, pthis.Biheader.SizeImage)
	err = binary.Read(file, binary.LittleEndian, pthis.BmpData)
	if err != nil {
		LogDebug(err)
		return err
	}
	return nil
}

func (pthis*Bitmap)WriteBmp(bmpFile string, binData[]uint8, w int, h int, bitCount int)error{
	file, err := os.Create(bmpFile)
	if err != nil {
		LogDebug(err)
		return err
	}

	defer func() {
		file.Sync()
		file.Close()
	}()

	// write BM
	err = binary.Write(file, binary.LittleEndian, uint8('B'))
	if err != nil {
		LogDebug(err)
		return err
	}
	err = binary.Write(file, binary.LittleEndian, uint8('M'))
	if err != nil {
		LogDebug(err)
		return err
	}

	headerLen := BITMAP_HEADER_TOTAL
	if pthis.BinPalette != nil {
		headerLen += len(pthis.BinPalette)
	}
	//write total file size
	err = binary.Write(file, binary.LittleEndian, uint32(headerLen + len(binData)))
	if err != nil {
		LogDebug(err)
		return err
	}

	//write reserver
	err = binary.Write(file, binary.LittleEndian, uint16(0))
	if err != nil {
		LogDebug(err)
		return err
	}
	err = binary.Write(file, binary.LittleEndian, uint16(0))
	if err != nil {
		LogDebug(err)
		return err
	}

	//write offset
	err = binary.Write(file, binary.LittleEndian, uint32(headerLen))
	if err != nil {
		LogDebug(err)
		return err
	}

	//write bitmap info
	err = binary.Write(file, binary.LittleEndian, &pthis.Biheader)
	if err != nil {
		LogDebug(err)
		return err
	}

	//write palette
	err = binary.Write(file, binary.LittleEndian, pthis.BinPalette)
	if err != nil {
		LogDebug(err)
		return err
	}

	//write bitmap bin data
	err = binary.Write(file, binary.LittleEndian, binData)
	if err != nil {
		LogDebug(err)
		return err
	}
	return nil
}
