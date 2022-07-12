package main

import (
	"encoding/json"
	"fmt"
	"lin/lin_common"
	"time"
)

func test_bmp(){

	bmp := &lin_common.Bitmap{}
	bmp.ReadBmp("../resource/aa.bmp")
	binJson,_ := json.Marshal(bmp.Biheader)
	fmt.Println(string(binJson))
	binJson,_ = json.Marshal(bmp.BfHeader)
	fmt.Println(string(binJson))
	fmt.Println(bmp.Biheader, bmp.BfHeader) // {40 384 290 1 24 0 334080 0 0 0 0} {19778 334134 0 0 54}

	bmp.WriteBmp("../resource/testw.bmp"/*, bmp.BmpData, int(bmp.Biheader.Width), int(bmp.Biheader.Height), int(bmp.Biheader.BitCount)*/)
}

func test_map(){
	testMap := &lin_common.MapData{}
	testMap.LoadMap("../resource/sample.bmp")

	bret := testMap.IsBlock(0, 0)
	fmt.Println(bret)

	bret = testMap.IsBlock(1, 0)
	fmt.Println(bret)

	bret = testMap.IsBlock(0, 1)
	fmt.Println(bret)
	bret = testMap.IsBlock(0, 2)
	fmt.Println(bret)

	bret = testMap.IsBlock(1, 1)
	fmt.Println(bret)

	testMap.DumpMap("../resource/dump.bmp", nil, nil, nil, nil)

	//src := lin_common.Coord2d{10, 290 - 261}
	//dst := lin_common.Coord2d{367,290 - 109}
	src := lin_common.Coord2d{72, 342 - 158}
	dst := lin_common.Coord2d{252,342 - 157}
	{
		fmt.Println("begin search")
		t1 := time.Now().UnixMilli()
		path, searchMgr := testMap.PathSearch(src, dst)
		t2 := time.Now().UnixMilli()
		fmt.Println("end search:", t2 - t1)
/*		for _, val := range path {
			fmt.Println(val)
		}
*/		testMap.DumpMap("../resource/path.bmp", path, &src, &dst, searchMgr)
	}

	{
		path := testMap.PathJPS(src, dst)
		fmt.Println("end search")
		for _, val := range path {
			fmt.Println(val)
		}
		testMap.DumpMap("../resource/jumppath.bmp", path, &src, &dst, nil)
	}

}
