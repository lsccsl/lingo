package main

import (
	"goserver/common"
	"goserver/msgpacket"
	"reflect"
)

func (pthis*Server)processOtherServerMsg (interMsg *interProtoMsg){
	switch t:=interMsg.protoMsg.(type){
	case *msgpacket.MSG_TEST_RPC:
		pthis.processTestRPC(t)
	}
}

func (pthis*Server)processTestRPC(msg *msgpacket.MSG_TEST_RPC){
	common.LogDebug("test rpc, srv:", pthis.srvID, " count:", msg.RpcCount)

	if pthis.connDial == nil {
		return
	}

	oldConnID := pthis.connDial.TcpConnectionID()
	go func() {
		for i := 0; i < int(msg.RpcCount); i ++ {

			if pthis.connDial == nil {
				return
			}

			if oldConnID != pthis.connDial.TcpConnectionID() {
				return
			}

			uuid := common.GenGUID()
			msgRes, err := srvMgr.SendRPC_Async(pthis.srvID,
				msgpacket.MSG_TYPE__MSG_TEST,
				&msgpacket.MSG_TEST{Id:uuid,Seq:int64(i)},
				60 * 1000)
			if err != nil {
				common.LogDebug("rpc err, uuid:", uuid, " srv:", pthis.srvID, " err:", err)
				continue
			}
			switch t := msgRes.(type) {
			case *msgpacket.MSG_TEST_RES:
				if t.Id != uuid {
					common.LogErr("id err:", t.Id, " uuid:", uuid, " srv:", pthis.srvID)
				}
			default:
				{
					var typeName string
					typeT := reflect.TypeOf(t)
					if typeT != nil {
						typeName = typeT.String()
					}
					common.LogErr("rsp msg type err:", t, typeName)
				}
			}
		}
	}()
}