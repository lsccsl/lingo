package main

import (
	"lin/lin_common"
	"lin/msgpacket"
)

func (pthis*Server)processOtherServerMsg (interMsg * interProtoMsg){
	switch t:=interMsg.protoMsg.(type){
	case *msgpacket.MSG_TEST_RPC:
		pthis.processTestRPC(t)
	}
}

func (pthis*Server)processTestRPC(msg *msgpacket.MSG_TEST_RPC){
	lin_common.LogDebug("test rpc, srvid:", pthis.srvID)
	go func() {
		for i := 0; i < int(msg.RpcCount); i ++ {
			msgRes := srvMgr.SendRPC_Async(pthis.srvID,
				msgpacket.MSG_TYPE__MSG_TEST,
				&msgpacket.MSG_TEST{Id:lin_common.GenGUID(),Seq:int64(i)},
				10 * 1000)
			lin_common.LogDebug(msgRes)
		}
	}()
}