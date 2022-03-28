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
	lin_common.LogDebug("test rpc, srvid:", pthis.srvID, " count:", msg.RpcCount)
	go func() {
		for i := 0; i < int(msg.RpcCount); i ++ {
			uuid := lin_common.GenGUID()
			msgRes := srvMgr.SendRPC_Async(pthis.srvID,
				msgpacket.MSG_TYPE__MSG_TEST,
				&msgpacket.MSG_TEST{Id:uuid,Seq:int64(i)},
				10 * 1000)
			switch t := msgRes.(type) {
			case *msgpacket.MSG_TEST_RES:
				if t.Id != uuid {
					lin_common.LogErr("id err:", t.Id, " uuid:", uuid, " srvid:", pthis.srvID)
				}
			default:
				lin_common.LogErr("rsp msg type err:", t)
			}
		}
	}()
}