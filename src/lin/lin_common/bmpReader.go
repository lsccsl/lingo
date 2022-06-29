package lin_common

import (
	"encoding/binary"
	"fmt"
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

type BitmapFileHeader struct {
	bfType uint16
	bfSize uint32
	bfReserved1 uint16
	bfReserved2 uint16
	bfOffBits uint32
}

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
}

type Bitmap struct{
	bfHeader BitmapFileHeader
	biheader BitmapInfoHeader
}

func (pthis*Bitmap)ReadBmp(bmpFile string){
	file, err := os.Open(bmpFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	//type拆成两个byte来读
	//Read第二个参数字节序一般windows/linux大部分都是LittleEndian,苹果系统用BigEndian
	//binary.Read(file, binary.LittleEndian, &pthis.headA)
	//binary.Read(file, binary.LittleEndian, &pthis.headB)
	binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfType)

	//文件大小
	binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfSize)

	//预留字节
	binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfReserved1)
	binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfReserved2)

	//偏移字节
	binary.Read(file, binary.LittleEndian, &pthis.bfHeader.bfOffBits)

	fmt.Println(pthis.bfHeader)

	binary.Read(file, binary.LittleEndian, &pthis.biheader)
	fmt.Println(pthis.biheader)
}