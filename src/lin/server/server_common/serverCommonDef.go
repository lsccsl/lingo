package server_common

import (
	"lin/lin_common"
	"strconv"
)

type MSGQUE_SRV_ID int64
const MSGQUE_SRV_ID_INVALID MSGQUE_SRV_ID = 0

const(
	EN_TCP_CLOSE_REASON_reg_reconnect = lin_common.EN_TCP_CLOSE_REASON_inter_max + 1
	EN_TCP_CLOSE_REASON_recv_ntf_offline = lin_common.EN_TCP_CLOSE_REASON_inter_max + 2
	EN_TCP_CLOSE_REASON_repeated_msgque_center = lin_common.EN_TCP_CLOSE_REASON_inter_max + 3
)

func (id MSGQUE_SRV_ID)String() string {
	return "[que_srv_id:" + strconv.FormatInt(int64(id), 10) + "]"
}