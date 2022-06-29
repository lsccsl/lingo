package main

import (
	"lin/lin_common"
	//"image"
)

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

func test_bmp(bmp_file string){

	bmp := &lin_common.Bitmap{}
	bmp.ReadBmp(bmp_file)

}