package lin_common

import (
	"encoding/hex"
	uuid_v1 "github.com/google/uuid"
	uuid_vX "github.com/satori/go.uuid"
	tsgutils "github.com/typa01/go-utils"
)

const UUID_INVALID int64 = 0

func GenUUID64_V1() int64 {
	newUUID, err := uuid_v1.NewUUID()
	if err != nil {
		return UUID_INVALID
	}

	if len(newUUID) < 16 {
		return UUID_INVALID
	}

	var uuidFirst int64 = 0
	{
		idx := 0
		for i := 3; i >= 0; i-- {
			var tmp int64
			tmp = int64(newUUID[i])
			tmp = tmp << (idx * 8)
			idx++
			uuidFirst = uuidFirst | tmp
		}
		for i := 5; i >= 4; i-- {
			var tmp int64
			tmp = int64(newUUID[i])
			tmp = tmp << (idx * 8)
			idx++
			uuidFirst = uuidFirst | tmp
		}
		for i := 7; i >= 6; i-- {
			var tmp int64
			tmp = int64(newUUID[i])
			tmp = tmp << (idx * 8)
			idx++
			uuidFirst = uuidFirst | tmp
		}
	}

	var uuidSecond int64 = 0
	{
		idx2 := 0
		for i := 9; i >= 8; i-- {
			var tmp int64
			tmp = int64(newUUID[i])
			tmp = tmp << (idx2 * 8)
			idx2++
			uuidSecond = uuidSecond | tmp
		}

		for i := 15; i >= 10; i-- {
			var tmp int64
			tmp = int64(newUUID[i])
			tmp = tmp << (idx2 * 8)
			idx2++
			uuidSecond = uuidSecond | tmp
		}
	}

	//fmt.Println(uuidFirst, uuidSecond)
	ret := (uuidFirst + uuidSecond)
	return ret
}

func GenUUID64_V4() int64 {
	newUUID := uuid_vX.NewV4()

	if len(newUUID) < 16 {
		return UUID_INVALID
	}

	var uuidFirst int64 = 0
	for i := 0; i < 8; i++ {
		var tmp int64
		tmp = int64(newUUID[i])
		tmp = tmp << (i * 8)
		uuidFirst = uuidFirst | tmp
	}
	var uuidSecond int64 = 0
	for i := 8; i < 16; i++ {
		var tmp int64
		tmp = int64(newUUID[i])
		tmp = tmp << ((i - 8) * 8)
		uuidSecond = uuidSecond | tmp
	}

	ret := uuidFirst ^ uuidSecond
	return ret
}

func GenGUID() int64 {
	sGuid := tsgutils.GUID()
	bin, err := hex.DecodeString(sGuid)
	if err != nil || bin == nil {
		return UUID_INVALID
	}
	if len(bin) < 16 {
		return UUID_INVALID
	}
	var uuidFirst int64 = 0
	for i := 0; i < 8; i++ {
		var tmp int64
		tmp = int64(bin[i])
		tmp = tmp << (i * 8)
		uuidFirst = uuidFirst | tmp
	}
	var uuidSecond int64 = 0
	for i := 8; i < 16; i++ {
		var tmp int64
		tmp = int64(bin[i])
		tmp = tmp << ((i - 8) * 8)
		uuidSecond = uuidSecond | tmp
	}

	ret := uuidFirst ^ uuidSecond
	return ret
}

