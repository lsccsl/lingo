package msgque_client

import (
	"github.com/golang/protobuf/proto"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
	"time"
)



func (pthis*MgrQueClient)SendMsg(srvUUIDTo server_common.SRV_ID, srvType server_common.SRV_TYPE,
	msgType msgpacket.PB_MSG_TYPE, protoMsg proto.Message, timeoutMilliSec int) (res proto.Message, err error) {
	res = nil
	err = nil
	defer func() {
		err := recover()
		if err != nil {
			common.LogErr(err)
		}
	}()
	if pthis.msgMgr == nil{
		err = common.GenErr(common.ERR_sys, "no rpc mgr")
		return
	}

	msgID := server_common.MSG_ID(common.GenUUID64_V4())
	pmsg := &msgpacket.PB_MSG_INTER_MSG{
		MsgId:int64(msgID),
		MsgType:int32(msgType),
		SrvUuidFrom:int64(pthis.srvUUID),
		SrvUuidTo:int64(srvUUIDTo),
		SrvType:int32(srvType),
		MsgSeq:pthis.msgMgr.seq.Add(1),
		Timestamp:time.Now().UnixMilli(),
		TimeoutWait:int64(timeoutMilliSec),
	}

	common.LogDebug("from:", pthis.srvUUID, "to:", srvUUIDTo, msgID, srvType, " msg type:", pmsg.MsgType)

	pmsg.MsgBin, err = proto.Marshal(protoMsg)
	if err != nil {
		common.LogErr(err)
		err = common.GenErr(common.ERR_sys, "packet err")
		return
	}

	var rreq *MsgReq = nil
	if timeoutMilliSec > 0 {
		rreq = pthis.msgMgr.ClientSrvMsgMgrAddReq(msgID)
	}

	pthis.SendProtoMsg(pthis.fdQueSrv, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG, pmsg)

	if timeoutMilliSec > 0 && rreq != nil {
		select {
		case resCh := <-rreq.chNtf:
			if resCh != nil {
				if resCh.Res == msgpacket.PB_RESPONSE_CODE_PB_RESPONSE_CODE_OK {
					res, _ = resCh.PBMsg.(proto.Message)
				} else {
					err = common.GenErr(common.ERR_rpc_response_err, " res:", resCh.Res)
					return
				}
			}
		case <-time.After(time.Millisecond * time.Duration(timeoutMilliSec)):
			err = common.GenErr(common.ERR_rpc_timeout, " msg time out srv:", srvUUIDTo, msgID)
		}

		pthis.msgMgr.ClientSrvMsgMgrDelReq(msgID)
	}

	return
}


func (pthis*MgrQueClient)process_PB_MSG_INTER_MSG(fd common.FD_DEF, pbMsg proto.Message, inAttachData interface{}) {

	msgReq, ok := pbMsg.(*msgpacket.PB_MSG_INTER_MSG)
	if !ok || msgReq == nil {
		common.LogErr("msg convert err")
		return
	}
	msgReq.TimestampArrive = time.Now().UnixMilli()

	srvUUIDFrom := server_common.SRV_ID(msgReq.SrvUuidFrom)
	srvUUIDTo := server_common.SRV_ID(msgReq.SrvUuidTo)
	msgID := server_common.MSG_ID(msgReq.MsgId)

	common.LogInfo("from:", srvUUIDFrom, "to:", srvUUIDTo, msgID, " msg type:", msgReq.MsgType)
	protoMsg := msgpacket.ParseProtoMsg(msgReq.MsgBin, msgReq.MsgType)
	common.LogInfo("packType:", msgReq.MsgType, " protoMsg:", protoMsg)

	pthis.Go_ProcessMsg(msgReq, msgReq.MsgType, protoMsg)
}

func (pthis*MgrQueClient)Go_ProcessMsg(msgReq *msgpacket.PB_MSG_INTER_MSG, pbMsgType int32, msgBody proto.Message) {
	go func() {

		defer func() {
			err := recover()
			if err != nil {
				common.LogErr(err)
			}
		}()

		if pthis.cb != nil {
			if msgReq.TimeoutWait > 0 {
				msgRes := &msgpacket.PB_MSG_INTER_MSG_RES{
					SrvUuidFrom:     msgReq.SrvUuidFrom,
					SrvUuidTo:       msgReq.SrvUuidTo,
					SrvType:         msgReq.SrvType,
					MsgId:           msgReq.MsgId,
					MsgSeq:          msgReq.MsgSeq,
					Timestamp:       msgReq.Timestamp,
					TimestampArrive: msgReq.TimestampArrive,
					TimeoutWait:     msgReq.TimeoutWait,
					Res:             msgpacket.PB_RESPONSE_CODE_PB_RESPONSE_CODE_OK,
				}
				var msgBodyRes proto.Message
				msgRes.MsgType, msgBodyRes = pthis.cb.Go_CallBackMsg(msgBody, pbMsgType,
					server_common.SRV_ID(msgReq.SrvUuidFrom), server_common.SRV_TYPE(msgReq.SrvType), int(msgReq.TimeoutWait))

				if msgRes != nil {
					var err error
					msgRes.MsgBin, err = proto.Marshal(msgBodyRes)
					if err != nil {
						common.LogErr(err)
					}
				}
				pthis.SendProtoMsg(pthis.fdQueSrv, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG_RES, msgRes)
			} else {
				pthis.cb.Go_CallBackMsg(msgBody, pbMsgType,
					server_common.SRV_ID(msgReq.SrvUuidFrom), server_common.SRV_TYPE(msgReq.SrvType), 0)
			}
		} else {
			common.LogDebug("no call back")
		}
	}()
}


func (pthis*MgrQueClient)process_PB_MSG_INTER_MSG_RES(fd common.FD_DEF, pbMsg proto.Message, inAttachData interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			common.LogErr(err)
		}
	}()

	msgRes, ok := pbMsg.(*msgpacket.PB_MSG_INTER_MSG_RES)
	if !ok || msgRes == nil {
		common.LogErr("msg convert err")
		return
	}

	srvUUIDFrom := server_common.SRV_ID(msgRes.SrvUuidFrom)
	srvUUIDTo := server_common.SRV_ID(msgRes.SrvUuidTo)
	msgID := server_common.MSG_ID(msgRes.MsgId)

	common.LogInfo("from:", srvUUIDFrom, "to:", srvUUIDTo, msgID, " msg type:", msgRes.MsgType)
	protoMsg := msgpacket.ParseProtoMsg(msgRes.MsgBin, msgRes.MsgType)
	common.LogInfo("packType:", msgRes.MsgType, " protoMsg:", protoMsg)

	if pthis.msgMgr == nil {
		return
	}
	rreq := pthis.msgMgr.ClientSrvMsgMgrFindReq(msgID)
	if rreq == nil {
		common.LogErr("fail find rpc:", msgID, " srv:", srvUUIDTo)
		return
	}
	if rreq.chNtf != nil {
		rreq.chNtf <- &MsgRes{Res: msgRes.Res, PBMsg:protoMsg}
	}
}
