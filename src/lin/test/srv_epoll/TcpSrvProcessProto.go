package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"reflect"
)

func (pthis*TcpSrv)process_ProtoMsg(fd lin_common.FD_DEF, protoMsg proto.Message) {
	switch t:=protoMsg.(type){
	case *msgpacket.MSG_TEST_RPC:
		pthis.process_MSG_TEST_RPC(fd, t)
	}
}


func (pthis*TcpSrv)process_MSG_TEST_RPC(fd lin_common.FD_DEF, msg *msgpacket.MSG_TEST_RPC) {
	lin_common.LogDebug("test rpc, srv:", pthis.srvID, " count:", msg.RpcCount)

	if pthis.fdDial.IsNull() {
		lin_common.LogDebug(" tcp close, srv:", pthis.srvID)
		return
	}

	oldConnID := pthis.fdDial
	seq := 0
	go func() {
		for /*i := 0; i < int(msg.RpcCount); i ++*/ {
			if pthis.fdDial.IsNull() {
				lin_common.LogDebug(" tcp close, srv:", pthis.srvID)
				return
			}
			if !pthis.fdDial.IsSame(&oldConnID) {
				lin_common.LogDebug(" tcp close, srv:", pthis.srvID)
				return
			}

			seq++
			uuid := lin_common.GenGUID()
			//lin_common.LogDebug(" send rpc :", uuid)
			msgRes, err := pthis.pu.tcpSrvMgr.TcpSrvMgrRPCSync(pthis.srvID,
				msgpacket.MSG_TYPE__MSG_TEST,
				&msgpacket.MSG_TEST{Id:uuid,Seq:int64(seq)},
				60 * 1000)
			//.LogDebug(" end rpc :", uuid)
			if err != nil {
				lin_common.LogDebug("rpc err, uuid:", uuid, " srv:", pthis.srvID, " err:", err)
				continue
			}
			switch t := msgRes.(type) {
			case *msgpacket.MSG_TEST_RES:
				if t.Id != uuid {
					lin_common.LogErr("id err:", t.Id, " uuid:", uuid, " srv:", pthis.srvID)
				} else {
					//lin_common.LogDebug(" recv rpc res:", t)
				}
			default:
				{
					var typeName string
					typeT := reflect.TypeOf(t)
					if typeT != nil {
						typeName = typeT.String()
					}
					lin_common.LogErr("rsp msg type err:", t, typeName)
				}
			}
		}
	}()
}