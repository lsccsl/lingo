package server_common

import (
	"lin/lin_common"
	"strconv"
)

type SRV_ID int64
const SRV_ID_INVALID SRV_ID = 0

type SRV_TYPE int32

const(
	EN_TCP_CLOSE_REASON_reg_reconnect = lin_common.EN_TCP_CLOSE_REASON_inter_max + 1
	EN_TCP_CLOSE_REASON_recv_ntf_offline = lin_common.EN_TCP_CLOSE_REASON_inter_max + 2
	EN_TCP_CLOSE_REASON_repeated_msgque_center = lin_common.EN_TCP_CLOSE_REASON_inter_max + 3
)

func (id SRV_ID)String() string {
	return "[srv_id:" + strconv.FormatInt(int64(id), 10) + "]"
}

func (st SRV_TYPE)String() string {
	return "[srv_type:" + strconv.FormatInt(int64(st), 10) + "]"
}

