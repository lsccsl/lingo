package main

import (
	"goserver/server/server_common"
	msgque_client "goserver/server/server_msg_que_client"
)

type DBSrv struct {
	mqClient *msgque_client.MgrQueClient

	dbMgr * DBMongoMgr

	dbDataTypeMgr *DBDataTypeMgr
}

func (pthis*DBSrv)Wait() {
	pthis.mqClient.WaitEpoll()
}



func ConstructDBSrv(id string)*DBSrv {

	cs := &DBSrv{
		dbMgr:ConstructorDBMongoMgr(),
		dbDataTypeMgr:ConstructDBDataTypeMgr(),
	}

	cs.mqClient = msgque_client.ConstructMgrQueClient(server_common.Global_ServerCfg.MsgQueCent.OutAddr, server_common.SRV_TYPE_database_server, cs)
	cs.mqClient.DialToQueSrv()

	return cs
}

func (pthis*DBSrv)Dump(bDetail bool) string {
	str := pthis.mqClient.Dump(bDetail)
	return str
}
