package server_common

import "lin/lin_common"

type MSGQUE_SRV_ID int64
const MSGQUE_SRV_ID_INVALID MSGQUE_SRV_ID = 0

const(
	EN_TCP_CLOSE_REASON_reg_reconnect = lin_common.EN_TCP_CLOSE_REASON_inter_max + 1
	EN_TCP_CLOSE_REASON_recv_ntf_offline = lin_common.EN_TCP_CLOSE_REASON_inter_max + 2
)