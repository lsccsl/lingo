package main

import (
	"fmt"
	"lin/lin_common"
	"runtime"
	"time"
)

func test_uuid(){
	test_uuid1()
	test_uuid2()
}

func test_uuid2() {
	fmt.Println("uuid v1")
	for i := 0; i < 10 ; i ++ {
		uuid := lin_common.GenUUID64_V1()
		fmt.Println(uuid)
	}

	fmt.Println("uuid v4")
	for i := 0; i < 10 ; i ++ {
		uuid := lin_common.GenUUID64_V4()
		fmt.Println(uuid)
	}

	fmt.Println("guid")
	for i := 0; i < 10 ; i ++ {
		uuid := lin_common.GenGUID()
		fmt.Println(uuid)
	}
}

func test_uuid1() {
	/*{
		fmt.Println("uuid v1")
		for j := 0; j < 1000; j ++ {
			runtime.GC()
			mapUUID := make(map[int64]int64)
			count := 0
			for i := 0; i < 1000 * 1000; i ++ {
				count ++
				uuid := lin_common.GenUUID64_V1()
				mapUUID[uuid] = uuid
			}

			fmt.Print(".")
			if len(mapUUID) != count {
				fmt.Println("uuid v1 repeated:", count, " len:", len(mapUUID))
				time.Sleep(time.Second * 5)
				break
			}
		}
	}*/

	{
		fmt.Println("uuid v4")
		for j := 0; j < 1000; j ++ {
			runtime.GC()
			count := 0
			mapUUID := make(map[int64]int64)
			for i := 0; i < 1000 * 1000; i ++ {
				count ++
				uuid := lin_common.GenGUID()
				mapUUID[uuid] = uuid
			}
			fmt.Print(".")

			if len(mapUUID) != count {
				fmt.Println("uuid 4 repeated:", count, " len:", len(mapUUID))
				time.Sleep(time.Second * 5)
				break
			}
		}
	}

	{
		fmt.Println("guid")
		for j := 0; j < 1000; j ++ {
			runtime.GC()
			count := 0
			mapGUID := make(map[int64]int64)
			for i := 0; i < 1000 * 1000; i ++ {
				count ++
				uuid := lin_common.GenUUID64_V4()
				mapGUID[uuid] = uuid
			}
			fmt.Print(".")

			if len(mapGUID) != count {
				fmt.Println("uuid repeated:", count, " len:", len(mapGUID))
				time.Sleep(time.Second * 5)
				break
			}
		}
	}
}
