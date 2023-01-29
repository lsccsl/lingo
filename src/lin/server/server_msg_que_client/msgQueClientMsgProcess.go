package msgque_client

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"time"
)


func (pthis*MgrQueClient)SendMsg(srvUUIDTo server_common.SRV_ID, srvType server_common.SRV_TYPE,
	msgType msgpacket.PB_MSG_TYPE, protoMsg proto.Message, timeoutMilliSec int) (proto.Message, error) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()
	if pthis.msgMgr == nil{
		return nil, lin_common.GenErr(lin_common.ERR_sys, "no rpc mgr")
	}

	msgID := server_common.MSG_ID(lin_common.GenUUID64_V4())
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
	var err error
	pmsg.MsgBin, err = proto.Marshal(protoMsg)
	if err != nil {
		lin_common.LogErr(err)
		return nil, lin_common.GenErr(lin_common.ERR_sys, "packet err")
	}

	var rreq *MsgReq = nil
	if timeoutMilliSec > 0 {
		rreq = pthis.msgMgr.ClientSrvMsgMgrAddReq(msgID)
	}

	pthis.SendProtoMsg(pthis.fdQueSrv, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG, pmsg)

	if timeoutMilliSec > 0 && rreq != nil {

		var res proto.Message = nil
		select {
		case resCh := <-rreq.chNtf:
			res, _ = resCh.(proto.Message)
		case <-time.After(time.Millisecond * time.Duration(timeoutMilliSec)):
			err = lin_common.GenErr(lin_common.ERR_rpc_timeout, " msg time out srv:", srvUUIDTo, msgID)
		}

		pthis.msgMgr.ClientSrvMsgMgrDelReq(msgID)

		if err != nil {
			return nil, err
		}
		if res == nil {
			return nil, lin_common.GenErr(lin_common.ERR_sys, "msg is nil")
		}
		return res, nil
	}

	return nil, nil
}


func (pthis*MgrQueClient)process_PB_MSG_INTER_MSG(fd lin_common.FD_DEF, pbMsg proto.Message, inAttachData interface{}) {

	msgReq, ok := pbMsg.(*msgpacket.PB_MSG_INTER_MSG)
	if !ok || msgReq == nil {
		lin_common.LogErr("msg convert err")
		return
	}
	msgReq.TimestampArrive = time.Now().UnixMilli()

	srvUUIDFrom := server_common.SRV_ID(msgReq.SrvUuidFrom)
	srvUUIDTo := server_common.SRV_ID(msgReq.SrvUuidTo)
	msgID := server_common.MSG_ID(msgReq.MsgId)

	lin_common.LogInfo("from:", srvUUIDFrom, "to:", srvUUIDTo, msgID, " msg type:", msgReq.MsgType)
	protoMsg := msgpacket.ParseProtoMsg(msgReq.MsgBin, msgReq.MsgType)
	lin_common.LogInfo("packType:", msgReq.MsgType, " protoMsg:", protoMsg)

	pthis.Go_ProcessMsg(msgReq, msgReq.MsgType, protoMsg)
}

func (pthis*MgrQueClient)Go_ProcessMsg(msgReq *msgpacket.PB_MSG_INTER_MSG, pbMsgType int32, msgBody proto.Message) {
	go func() {

		defer func() {
			err := recover()
			if err != nil {
				lin_common.LogErr(err)
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
				msgRes.MsgType, msgBodyRes = pthis.cb.Go_CallBackProcessMsg(msgBody, pbMsgType,
					server_common.SRV_ID(msgReq.SrvUuidFrom), server_common.SRV_TYPE(msgReq.SrvType), int(msgReq.TimeoutWait))

				if msgRes != nil {
					var err error
					msgRes.MsgBin, err = proto.Marshal(msgBodyRes)
					if err != nil {
						lin_common.LogErr(err)
					}
				}
				pthis.SendProtoMsg(pthis.fdQueSrv, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG_RES, msgRes)
			} else {
				pthis.cb.Go_CallBackProcessMsg(msgBody, pbMsgType,
					server_common.SRV_ID(msgReq.SrvUuidFrom), server_common.SRV_TYPE(msgReq.SrvType), 0)
			}
		} else {
			lin_common.LogDebug("no call back")
		}
	}()
}


func (pthis*MgrQueClient)process_PB_MSG_INTER_MSG_RES(fd lin_common.FD_DEF, pbMsg proto.Message, inAttachData interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	msgReq, ok := pbMsg.(*msgpacket.PB_MSG_INTER_MSG_RES)
	if !ok || msgReq == nil {
		lin_common.LogErr("msg convert err")
		return
	}
	msgReq.TimestampArrive = time.Now().UnixMilli()

	srvUUIDFrom := server_common.SRV_ID(msgReq.SrvUuidFrom)
	srvUUIDTo := server_common.SRV_ID(msgReq.SrvUuidTo)
	msgID := server_common.MSG_ID(msgReq.MsgId)

	lin_common.LogInfo("from:", srvUUIDFrom, "to:", srvUUIDTo, msgID, " msg type:", msgReq.MsgType)
	protoMsg := msgpacket.ParseProtoMsg(msgReq.MsgBin, msgReq.MsgType)
	lin_common.LogInfo("packType:", msgReq.MsgType, " protoMsg:", protoMsg)

	if pthis.msgMgr == nil {
		return
	}
	rreq := pthis.msgMgr.ClientSrvMsgMgrFindReq(msgID)
	if rreq == nil {
		lin_common.LogErr("fail find rpc:", msgID, " srv:", srvUUIDTo)
		return
	}
	if rreq.chNtf != nil {
		rreq.chNtf <- protoMsg
	}
}
