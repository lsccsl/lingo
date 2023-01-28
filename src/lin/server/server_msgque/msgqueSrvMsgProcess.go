package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"time"
)

func (pthis*MsgQueSrv)process_PB_MSG_INTER_MSG(fd lin_common.FD_DEF, pbMsg proto.Message, inAttachData interface{}){
	lin_common.LogDebug(fd, " msg:", pbMsg, " attachData:", inAttachData)

	pmsg, ok := pbMsg.(*msgpacket.PB_MSG_INTER_MSG)
	if !ok || pmsg == nil {
		lin_common.LogErr("msg convert err")
		return
	}
	srvUUIDFrom := server_common.SRV_ID(pmsg.SrvUuidFrom)
	srvUUIDTo := server_common.SRV_ID(pmsg.SrvUuidTo)
	msgID := server_common.MSG_ID(pmsg.MsgId)
	srvType := server_common.SRV_TYPE(pmsg.SrvType)

	lin_common.LogInfo("from:", srvUUIDFrom, "to:", srvUUIDTo, msgID, srvType)

	if srvUUIDTo == server_common.SRV_ID_INVALID {
		if srvType == server_common.SRV_TYPE_msg_center {
			pthis.SendProtoMsg(pthis.fdCenter, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG, pmsg)
		} else {
			pthis.processMsgLocal(fd, pmsg, inAttachData)
		}
	} else {
		fdRoute, ok := pthis.smgr.findLocalRoute(srvUUIDTo)
		if ok {
			pthis.SendProtoMsg(fdRoute, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG, pmsg)
			return
		} else {
			queSrvID := pthis.smgr.findRemoteRoute(srvUUIDTo)
			if server_common.SRV_ID_INVALID == queSrvID {
				lin_common.LogErr("can't find srv uuid", srvUUIDTo, msgID)
				return
			}

			qsi := otherMsgQueSrvInfo{}
			ok = pthis.otherMgr.Load(queSrvID, &qsi)
			if !ok {
				lin_common.LogErr("can't find que srv", queSrvID, srvUUIDTo, msgID)
				return
			}
			pthis.SendProtoMsg(qsi.fdDial, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG, pmsg)
		}
	}
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_MSG_RES(fd lin_common.FD_DEF, pbMsg proto.Message, inAttachData interface{}) {

	pmsg, ok := pbMsg.(*msgpacket.PB_MSG_INTER_MSG_RES)
	if !ok || pmsg == nil {
		lin_common.LogErr("msg convert err")
		return
	}
	srvUUIDFrom := server_common.SRV_ID(pmsg.SrvUuidFrom)
	srvUUIDTo := server_common.SRV_ID(pmsg.SrvUuidTo)
	msgID := server_common.MSG_ID(pmsg.MsgId)

	lin_common.LogInfo("from:", srvUUIDFrom, "to:", srvUUIDTo, msgID)
	fdRoute, ok := pthis.smgr.findLocalRoute(srvUUIDFrom)
	if ok {
		pthis.SendProtoMsg(fdRoute, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG_RES, pmsg)
		return
	} else {
		queSrvID := pthis.smgr.findRemoteRoute(srvUUIDFrom)
		if server_common.SRV_ID_INVALID == queSrvID {
			lin_common.LogErr("can't find srv uuid", srvUUIDFrom, msgID)
			return
		}

		qsi := otherMsgQueSrvInfo{}
		ok = pthis.otherMgr.Load(queSrvID, &qsi)
		if !ok {
			lin_common.LogErr("can't find que srv", queSrvID, srvUUIDFrom, msgID)
			return
		}
		pthis.SendProtoMsg(qsi.fdDial, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG_RES, pmsg)
	}
}

func (pthis*MsgQueSrv)processMsgLocal(fd lin_common.FD_DEF, msgReq * msgpacket.PB_MSG_INTER_MSG, inAttachData interface{}) {
	msgReq.TimestampArrive = time.Now().UnixMilli()

	srvUUIDFrom := server_common.SRV_ID(msgReq.SrvUuidFrom)
	srvUUIDTo := server_common.SRV_ID(msgReq.SrvUuidTo)
	msgID := server_common.MSG_ID(msgReq.MsgId)

	lin_common.LogInfo("from:", srvUUIDFrom, "to:", srvUUIDTo, msgID, " msg type:", msgReq.MsgType)
	msgBody := msgpacket.ParseProtoMsg(msgReq.MsgBin, msgReq.MsgType)
	lin_common.LogInfo("packType:", msgReq.MsgType, " protoMsg:", msgBody)

	var msgType int32 = 0
	var msgBodyRes proto.Message = nil
	switch t := msgBody.(type) {
	case *msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE:
		{
			//time.Sleep(time.Second * 3)
			msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_INTER_QUESRV_GET_SRVTYPE_RES)
			msgBodyRes = pthis.processMsgLocal_PB_MSG_INTER_QUESRV_GET_SRVTYPE(t, inAttachData)
		}
	}

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
		MsgType:         msgType,
	}

	if msgBodyRes != nil {
		var err error
		msgRes.MsgBin, err = proto.Marshal(msgBodyRes)
		if err != nil {
			lin_common.LogErr(err)
		}
	}

	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG_RES, msgRes)
}


func (pthis*MsgQueSrv)processMsgLocal_PB_MSG_INTER_QUESRV_GET_SRVTYPE(pbMsg * msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE,
	inAttachData interface{}) proto.Message {

	res := &msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE_RES{}

	arraySrvID := pthis.smgr.findSrvByType(server_common.SRV_TYPE(pbMsg.SrvType))
	if arraySrvID == nil {
		return nil
	}
	if len(arraySrvID) == 0 {
		return nil
	}

	for _, v := range arraySrvID {
		res.ArrarySrv = append(res.ArrarySrv, &msgpacket.PB_SRV_INFO_ONE{
			SrvUuid : int64(v.srvUUID),
			SrvType : int32(v.srvType),
		})
	}

	return res
}