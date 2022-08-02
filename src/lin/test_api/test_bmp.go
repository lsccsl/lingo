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
	//testMap.LoadMap("../resource/sample.bmp")
	testMap.LoadMap("../resource/aa.bmp")

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
	//dst := lin_common.Coord2d{338,290 - 18}
	src := lin_common.Coord2d{10, 290 - 261}
	dst := lin_common.Coord2d{367,290 - 109}
	//src := lin_common.Coord2d{72, 342 - 158}
	//dst := lin_common.Coord2d{252,342 - 157}



	{
		t1 := time.Now().UnixMilli()
		path, jpsMgr := testMap.PathJPS(src, dst)
		t2 := time.Now().UnixMilli()
		fmt.Println("jps end search:", t2 - t1)
		var pathConn []lin_common.Coord2d
		for i := 0; i < len(path) - 1; i ++ {
			pos1 := path[i]
			pos2 := path[i + 1]
			posDiff := pos1.Dec(&pos2)
			if posDiff.X > 0 {
				posDiff.X = 1
			}
			if posDiff.X < 0 {
				posDiff.X = -1
			}
			if posDiff.Y > 0 {
				posDiff.Y = 1
			}
			if posDiff.Y < 0 {
				posDiff.Y = -1
			}
			curPos := pos2
			for {
				pathConn = append(pathConn, curPos)
				if curPos.IsNear(&pos1) {
					break
				}
				curPos = curPos.Add(&posDiff)
			}
		}

		testMap.DumpMap("../resource/jump_path.bmp", pathConn, &src, &dst, nil)
		testMap.DumpJPSMap("../resource/jump_tree_path.bmp", nil, jpsMgr)
	}

	{
		fmt.Println("begin search")
		t1 := time.Now().UnixMilli()
		path, searchMgr := testMap.PathSearchAStart(src, dst)
		t2 := time.Now().UnixMilli()
		fmt.Println("a* end search:", t2 - t1)
		testMap.DumpMap("../resource/path.bmp", path, &src, &dst, searchMgr)
	}
}


/*
2022-08-02T14:46:02.464645823+08:00[TcpClient.go:96] route:106 main.(*TcpClient).Process_MSG_PATH_SEARCH path searchpos_src:{pos_x:18 pos_y:277} pos_dst:{pos_x:63 pos_y:17}
ERROR 2022-08-02T14:46:02.465792417+08:00[MapJPS.go:102] route:106 lin/lin_common.(*JSPMgr).addNode already in history
goroutine 106 [running]:
runtime/debug.Stack()
        /home/lsc/go/go/src/runtime/debug/stack.go:24 +0x65
lin/lin_common.LogErr({0xc00014fb78, 0x1, 0x1})
        /home/lsc/go/lin/lingo/src/lin/lin_common/log.go:63 +0x159
lin/lin_common.(*JSPMgr).addNode(0xc0002224b0, 0xc00048b7a0, 0xc00048b740)
        /home/lsc/go/lin/lingo/src/lin/lin_common/MapJPS.go:102 +0x29f
lin/lin_common.(*MapData).searchHorVer(0xc00040def0, 0xc0002224b0, 0xc00048b6e0, {0xc000622270, 0xc00014fd68}, 0x0, 0x1)
        /home/lsc/go/lin/lingo/src/lin/lin_common/MapJPS.go:370 +0x405
lin/lin_common.(*MapData).PathJPS(0xc00040def0, {0x2, 0x2}, {0x3f, 0x11})
        /home/lsc/go/lin/lingo/src/lin/lin_common/MapJPS.go:562 +0xad3
main.(*TcpClient).Process_MSG_PATH_SEARCH(0xc000222140, 0xc0006200c0)
        /home/lsc/go/lin/lingo/src/lin/srv_epoll/TcpClient.go:100 +0x118
main.(*TcpClient).Process_protoMsg(0xc000222140, 0xc000620100)
        /home/lsc/go/lin/lingo/src/lin/srv_epoll/TcpClient.go:160 +0x14a
main.(*TcpClientMgrUnit)._go_Process_unit(0xc00042c990)
        /home/lsc/go/lin/lingo/src/lin/srv_epoll/TcpClientMgrUnit.go:120 +0xe5
created by main.ConstructorEpollServerMgr
        /home/lsc/go/lin/lingo/src/lin/srv_epoll/ServerMgr.go:280 +0x47f

2022-08-02T14:46:03.991751849+08:00[TcpClient.go:96] route:106 main.(*TcpClient).Process_MSG_PATH_SEARCH path searchpos_src:{pos_x:63 pos_y:17} pos_dst:{pos_x:129 pos_y:275}
2022-08-02T14:46:05.878446634+08:00[TcpClient.go:96] route:106 main.(*TcpClient).Process_MSG_PATH_SEARCH path searchpos_src:{pos_x:129 pos_y:275} pos_dst:{pos_x:53 pos_y:18}
ERROR 2022-08-02T14:46:05.879353639+08:00[MapJPS.go:102] route:106 lin/lin_common.(*JSPMgr).addNode already in history
goroutine 106 [running]:
runtime/debug.Stack()
        /home/lsc/go/go/src/runtime/debug/stack.go:24 +0x65
lin/lin_common.LogErr({0xc00015bb78, 0x1, 0x1})
        /home/lsc/go/lin/lingo/src/lin/lin_common/log.go:63 +0x159
lin/lin_common.(*JSPMgr).addNode(0xc000222410, 0xc00060f440, 0xc00060f3e0)
        /home/lsc/go/lin/lingo/src/lin/lin_common/MapJPS.go:102 +0x29f
lin/lin_common.(*MapData).searchHorVer(0xc00040def0, 0xc000222410, 0xc00060e6c0, {0xc00041c1a0, 0xc00015bd68}, 0x2, 0x1)
        /home/lsc/go/lin/lingo/src/lin/lin_common/MapJPS.go:370 +0x405
lin/lin_common.(*MapData).PathJPS(0xc00040def0, {0x2, 0x2}, {0x35, 0x12})
        /home/lsc/go/lin/lingo/src/lin/lin_common/MapJPS.go:562 +0xad3
main.(*TcpClient).Process_MSG_PATH_SEARCH(0xc000222140, 0xc000620040)
        /home/lsc/go/lin/lingo/src/lin/srv_epoll/TcpClient.go:100 +0x118
main.(*TcpClient).Process_protoMsg(0xc000222140, 0xc000620080)
        /home/lsc/go/lin/lingo/src/lin/srv_epoll/TcpClient.go:160 +0x14a
main.(*TcpClientMgrUnit)._go_Process_unit(0xc00042c990)
        /home/lsc/go/lin/lingo/src/lin/srv_epoll/TcpClientMgrUnit.go:120 +0xe5
created by main.ConstructorEpollServerMgr
        /home/lsc/go/lin/lingo/src/lin/srv_epoll/ServerMgr.go:280 +0x47f



        2022-08-02T15:18:39.871443753+08:00[TcpClient.go:96] route:99 main.(*TcpClient).Process_MSG_PATH_SEARCH path searchpos_src:{pos_x:317  pos_y:236}  pos_dst:{pos_x:103  pos_y:227}
*/