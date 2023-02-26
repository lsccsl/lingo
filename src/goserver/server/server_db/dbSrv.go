package main

import (
	"goserver/server/server_common"
	msgque_client "goserver/server/server_msg_que_client"
)

type DBSrv struct {
	mqClient *msgque_client.MgrQueClient

	db * DBMongo
}

func (pthis*DBSrv)Wait() {
	pthis.mqClient.WaitEpoll()
}



func ConstructDBSrv(id string)*DBSrv {

	dbCfg := server_common.GetDBSrvCfg(id)

	cs := &DBSrv{
		db:ConstructorDBMongo(dbCfg.DBUser, dbCfg.DBPwd, dbCfg.DBIp, dbCfg.DBPort),
	}

	cs.mqClient = msgque_client.ConstructMgrQueClient(server_common.Global_ServerCfg.MsgQueCent.OutAddr, server_common.SRV_TYPE_database_server, cs)
	cs.mqClient.DialToQueSrv()

	return cs
}

func (pthis*DBSrv)Dump(bDetail bool) string {
	str := pthis.mqClient.Dump(bDetail)
	return str
}
