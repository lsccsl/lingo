package server_common

import (
	"strconv"
)


type SRV_ID int64
const SRV_ID_INVALID SRV_ID = 0

func (id SRV_ID)String() string {
	return "[srv_id:" + strconv.FormatInt(int64(id), 10) + "]"
}


type SRV_TYPE int32
const(
	SRV_TYPE_none            SRV_TYPE = 0
	SRV_TYPE_msq_que         SRV_TYPE = 1
	SRV_TYPE_msg_center      SRV_TYPE = 2
	SRV_TYPE_center_server   SRV_TYPE = 3
	SRV_TYPE_game_server     SRV_TYPE = 4
	SRV_TYPE_database_server SRV_TYPE = 5
	SRV_TYPE_logon_server    SRV_TYPE = 6
)
func (st SRV_TYPE)String() string {
	typeName := "SRV_TYPE_none"
	switch st {
	case SRV_TYPE_msq_que:
		typeName = "msq_que"
	case SRV_TYPE_msg_center:
		typeName = "msg_center"
	case SRV_TYPE_center_server:
		typeName = "center_server"
	case SRV_TYPE_game_server:
		typeName = "game_server"
	case SRV_TYPE_database_server:
		typeName = "database_server"
	case SRV_TYPE_logon_server:
		typeName = "logon_server"
	}
	return "[srv_type:" + typeName + "(" + strconv.FormatInt(int64(st), 10) + ")]"
}


type MSG_ID int64
func (id MSG_ID)String() string {
	return "[msg_id:" + strconv.FormatInt(int64(id), 10) + "]"
}

type SrvBaseInfo struct {
	SrvUUID SRV_ID
	SrvType SRV_TYPE
}
