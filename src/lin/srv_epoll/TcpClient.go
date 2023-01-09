package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"lin/navmeshwrapper"
	"net"
	"runtime"
	"time"
)

type MAP_CLIENT_STATIC map[msgpacket.MSG_TYPE]int64

type TcpClient struct {
	fd lin_common.FD_DEF
	addr net.Addr

	clientID int64
	objID int

	timerConnClose * time.Timer
	durationClose time.Duration
	pu *TcpClientMgrUnit

	timeLastActive int64
}

func ConstructorTcpClient(pu *TcpClientMgrUnit, fd lin_common.FD_DEF, clientID int64) *TcpClient {
	tc := &TcpClient{
		fd : fd,
		pu : pu,
		clientID : clientID,
		durationClose : time.Second*time.Duration(pu.eSrvMgr.clientCloseTimeoutSec),
		addr : lin_common.TcpGetPeerName(fd.FD),
		timeLastActive: time.Now().Unix(),
	}
	runtime.SetFinalizer(tc, (*TcpClient).Destructor)
	tc.timerConnClose = time.AfterFunc(tc.durationClose,
		func(){
			tnow := time.Now().Unix()
			lin_common.LogDebug("timeout close clientid:", tc.clientID, " fd:", tc.fd.String(),
				" timeLastActive:", tc.timeLastActive, " tnow:", tnow, " diff:", (tnow - tc.timeLastActive))
			tc.pu.eSrvMgr.lsn.EPollListenerCloseTcp(tc.fd, EN_TCP_CLOSE_REASON_timeout)
/*			tc.timerConnClose.Stop()
			tc.timerConnClose = nil*/
		})

	return tc
}

func (pthis*TcpClient)Destructor() {
	lin_common.LogDebug(" clientid:", pthis.clientID, " fd:", pthis.fd.String())
	runtime.SetFinalizer(pthis, nil)
	if pthis.timerConnClose != nil {
		pthis.timerConnClose.Stop()
		pthis.timerConnClose = nil
	}
}

func (pthis*TcpClient)Process_MSG_TCP_STATIC(msg *msgpacket.MSG_TCP_STATIC) {
	lin_common.LogDebug(" seq:", msg.Seq)

	msgRes := &msgpacket.MSG_TCP_STATIC_RES{
		ByteRecv:0,
		ByteProc:0,
		ByteSend:0,
	}
	pthis.pu.eSrvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_TCP_STATIC_RES, msgRes)
}
func (pthis*TcpClient)Process_MSG_TEST(msg *msgpacket.MSG_TEST) {
	//lin_common.LogDebug("clientid:", pthis.clientID, " fd:", pthis.fd.String())
	msgRes := &msgpacket.MSG_TEST_RES{}
	msgRes.Id = msg.Id
	msgRes.Str = msg.Str
	msgRes.Seq = msg.Seq
	msgRes.Timestamp = msg.Timestamp
	msgRes.TimestampArrive = msg.TimestampArrive
	msgRes.TimestampProcess = time.Now().UnixMilli()

	lin_common.TMP_tcpWrite(pthis.fd, msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_TEST_RES, msgRes))
	//pthis.pu.eSrvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_TEST_RES, msgRes)
}
func (pthis*TcpClient)Process_MSG_HEARTBEAT(msg *msgpacket.MSG_HEARTBEAT) {
	lin_common.LogDebug("clientid:", pthis.clientID, " fd:", pthis.fd.String())
	msgRes := &msgpacket.MSG_HEARTBEAT_RES{}
	msgRes.Id = msg.Id
	pthis.pu.eSrvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgRes)
}

func (pthis*TcpClient)Process_MSG_GET_MAP(msg * msgpacket.MSG_GET_MAP){
	lin_common.LogDebug("get map")
	msgRes := &msgpacket.MSG_GET_MAP_RES{}
	pthis.pu.eSrvMgr.mapMgr.GetMapProtoMsg(msgRes)
	pthis.pu.eSrvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_GET_MAP_RES, msgRes)
}

func (pthis*TcpClient)Process_MSG_PATH_SEARCH(msg * msgpacket.MSG_PATH_SEARCH){
	lin_common.LogDebug("path search", msg)
	mapData := pthis.pu.eSrvMgr.mapMgr.mapData
	src := lin_common.Coord2d{int(msg.PosSrc.PosX), int(msg.PosSrc.PosY)}
	dst := lin_common.Coord2d{int(msg.PosDst.PosX), int(msg.PosDst.PosY)}
	path, jpsMgr := mapData.PathJPS(src, dst)

	msgRes := &msgpacket.MSG_PATH_SEARCH_RES{}
	msgRes.PosSrc = msg.PosSrc
	msgRes.PosDst = msg.PosDst

	var pathConn []lin_common.Coord2d
	for i := len(path) - 1; i > 0; i -- {
		pos1 := path[i - 1]
		pos2 := path[i]

		msgRes.PathKeyPos = append(msgRes.PathKeyPos, &msgpacket.POS_T{PosX:int32(pos2.X), PosY:int32(pos2.Y)})

		posDiff := pos1.Dec(&pos2)
		if posDiff.X > 0 {
			posDiff.X = 1
		}
		if posDiff.X < 0 {
			posDiff.X = -1
		}
		if posDiff.Y > 0 {
			posDiff.Y = 1
		}
		if posDiff.Y < 0 {
			posDiff.Y = -1
		}
		curPos := pos2
		for {
			pathConn = append(pathConn, curPos)
			msgRes.PathPos = append(msgRes.PathPos, &msgpacket.POS_T{PosX:int32(curPos.X), PosY:int32(curPos.Y)})
			if curPos.IsNear(&pos1) {
				break
			}
			curPos = curPos.Add(&posDiff)
		}
	}
	msgRes.PathKeyPos = append(msgRes.PathKeyPos, &msgpacket.POS_T{PosX:int32(path[0].X), PosY:int32(path[0].Y)})

	mapData.DumpMap("../resource/Process_MSG_PATH_SEARCH.bmp", pathConn, &src, &dst, nil)
	mapData.DumpJPSMap("../resource/Process_MSG_PATH_SEARCH_tree.bmp", nil, jpsMgr)

	pthis.pu.eSrvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_PATH_SEARCH_RES, msgRes)
}

func (pthis*TcpClient)Process_MSG_NAV_SEARCH(msg *msgpacket.MSG_NAV_SEARCH) {
	lin_common.LogDebug("nav search", msg)

	msgSearch := &msgNavPathSearch{path:nil}
	if msg.PosSrc != nil {
		msgSearch.src = navmeshwrapper.Coord3f{msg.PosSrc.X, msg.PosSrc.Y, msg.PosSrc.Z}
	}
	if msg.PosDst != nil {
		msgSearch.dst = navmeshwrapper.Coord3f{msg.PosDst.X, msg.PosDst.Y, msg.PosDst.Z}
	}
	pthis.pu.eSrvMgr.mapProcMgr.addMapProcessMsg(msgSearch, pthis.clientID, time.Second * 3)

	lin_common.LogDebug(msgSearch.path)

	msg_ret := &msgpacket.MSG_NAV_SEARCH_RES{}
	if msgSearch.path != nil {
		for _, val := range msgSearch.path {
			msg_ret.PathPos = append(msg_ret.PathPos, &msgpacket.PROTO_VEC_3F{X:val.X, Y:val.Y, Z:val.Z})
		}
	}

	pthis.pu.eSrvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_NAV_SEARCH_RES, msg_ret)
}

func (pthis*TcpClient)Process_MSG_NAV_ADD_OBSTACLE(msg * msgpacket.MSG_NAV_ADD_OBSTACLE) {
	lin_common.LogDebug("add obstacle", msg)

	msgAdd := &msgNavAddObstacle{
		ob : NavObstacle{
			center :  navmeshwrapper.Coord3f{msg.Obstacle.Center.X,msg.Obstacle.Center.Y, msg.Obstacle.Center.Z},
			halfExt : navmeshwrapper.Coord3f{msg.Obstacle.HalfExt.X,msg.Obstacle.HalfExt.Y, msg.Obstacle.HalfExt.Z},
			yRadian : msg.Obstacle.YRadian,
		},
	}
	pthis.pu.eSrvMgr.mapProcMgr.addMapProcessMsg(msgAdd, pthis.clientID, time.Second * 3)

	obstacle_id := msgAdd.ob.obstacleID

	msg_ret := &msgpacket.MSG_NAV_ADD_OBSTACLE_RES{}
	msg_ret.Obstacle = &msgpacket.NAV_OBSTACLE{}
	msg_ret.Obstacle.ObstacleId = obstacle_id
	msg_ret.Obstacle.Center = msg.Obstacle.Center
	msg_ret.Obstacle.HalfExt = msg.Obstacle.HalfExt
	msg_ret.Obstacle.YRadian = msg.Obstacle.YRadian

	pthis.pu.eSrvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_NAV_ADD_OBSTACLE_RES, msg_ret)
}

func (pthis*TcpClient)Process_MSG_NAV_DEL_OBSTACLE(msg * msgpacket.MSG_NAV_DEL_OBSTACLE) {
	lin_common.LogDebug("del obstacle", msg)

	msgDel := &msgNavDelObstacle{obstacleID:msg.ObstacleId}
	pthis.pu.eSrvMgr.mapProcMgr.addMapProcessMsg(msgDel, pthis.clientID, time.Second * 3)

	msg_ret := &msgpacket.MSG_NAV_DEL_OBSTACLE_RES{}
	msg_ret.ObstacleId = msg.ObstacleId
	pthis.pu.eSrvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_NAV_DEL_OBSTACLE_RES, msg_ret)
}

func (pthis*TcpClient)Process_MSG_NAV_GET_ALL_OBSTACLE(msg * msgpacket.MSG_NAV_GET_ALL_OBSTACLE){
	lin_common.LogDebug("get obstacle", msg)

	msgGet := &msgNavGetAllObstacle{}
	pthis.pu.eSrvMgr.mapProcMgr.addMapProcessMsg(msgGet, pthis.clientID, time.Second * 3)

	msg_ret := &msgpacket.MSG_NAV_GET_ALL_OBSTACLE_RES{}
	for _,v := range msgGet.ob {
		ob := &msgpacket.NAV_OBSTACLE{}
		ob.ObstacleId = v.obstacleID
		ob.Center = &msgpacket.PROTO_VEC_3F{X:v.center.X, Y:v.center.Y, Z:v.center.Z}
		ob.HalfExt = &msgpacket.PROTO_VEC_3F{X:v.halfExt.X, Y:v.halfExt.Y, Z:v.halfExt.Z}
		ob.YRadian = v.yRadian
		msg_ret.Obstacle = append(msg_ret.Obstacle, ob)
	}
	pthis.pu.eSrvMgr.SendProtoMsg(pthis.fd, msgpacket.MSG_TYPE__MSG_NAV_GET_ALL_OBSTACLE_RES, msg_ret)
}

func (pthis*TcpClient)Process_MSG_UPDATE_POS(msg * msgpacket.MSG_UPDATE_POS){
	// todo
}

func (pthis*TcpClient)Process_protoMsg(msg *msgClient) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(" clientid:", pthis.clientID, " fd:", pthis.fd.String(), " err:", err, " msg:", msg)
		}
	}()
	pthis.timerConnClose.Reset(pthis.durationClose)
	pthis.timeLastActive = time.Now().Unix()

	switch t := msg.msg.(type) {
	case *msgpacket.MSG_TEST:
		pthis.Process_MSG_TEST(t)
	case *msgpacket.MSG_HEARTBEAT:
		pthis.Process_MSG_HEARTBEAT(t)
	case *msgpacket.MSG_TCP_STATIC:
		pthis.Process_MSG_TCP_STATIC(t)
	case *msgpacket.MSG_GET_MAP:
		pthis.Process_MSG_GET_MAP(t)
	case *msgpacket.MSG_PATH_SEARCH:
		pthis.Process_MSG_PATH_SEARCH(t)
	case *msgpacket.MSG_NAV_SEARCH:
		pthis.Process_MSG_NAV_SEARCH(t)
	case *msgpacket.MSG_NAV_ADD_OBSTACLE:
		pthis.Process_MSG_NAV_ADD_OBSTACLE(t)
	case *msgpacket.MSG_NAV_DEL_OBSTACLE:
		pthis.Process_MSG_NAV_DEL_OBSTACLE(t)
	case *msgpacket.MSG_NAV_GET_ALL_OBSTACLE:
		pthis.Process_MSG_NAV_GET_ALL_OBSTACLE(t)
	case *msgpacket.MSG_UPDATE_POS:
		pthis.Process_MSG_UPDATE_POS(t)
	}
}



