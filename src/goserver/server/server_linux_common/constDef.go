package server_linux_common

import "goserver/common"

const(
	EN_TCP_CLOSE_REASON_reg_reconnect = common.EN_TCP_CLOSE_REASON_inter_max + 1
	EN_TCP_CLOSE_REASON_recv_ntf_offline = common.EN_TCP_CLOSE_REASON_inter_max + 2
	EN_TCP_CLOSE_REASON_repeated_msgque_center = common.EN_TCP_CLOSE_REASON_inter_max + 3
	EN_TCP_CLOSE_REASON_msgque_center_ntf_offline = common.EN_TCP_CLOSE_REASON_inter_max + 4
	EN_TCP_CLOSE_REASON_srv_reg_ok = common.EN_TCP_CLOSE_REASON_inter_max + 5
)
