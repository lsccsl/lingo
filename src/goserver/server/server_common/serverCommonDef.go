package server_common

import (
	"goserver/common"
	"strconv"
)


type SRV_ID int64
const SRV_ID_INVALID SRV_ID = 0
const(
	EN_TCP_CLOSE_REASON_reg_reconnect = common.EN_TCP_CLOSE_REASON_inter_max + 1
	EN_TCP_CLOSE_REASON_recv_ntf_offline = common.EN_TCP_CLOSE_REASON_inter_max + 2
	EN_TCP_CLOSE_REASON_repeated_msgque_center = common.EN_TCP_CLOSE_REASON_inter_max + 3
	EN_TCP_CLOSE_REASON_msgque_center_ntf_offline = common.EN_TCP_CLOSE_REASON_inter_max + 4
	EN_TCP_CLOSE_REASON_srv_reg_ok = common.EN_TCP_CLOSE_REASON_inter_max + 5
)
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
