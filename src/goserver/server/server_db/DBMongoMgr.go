package main

import (
	"goserver/server/server_common"
	"sync"
)

type DBMongoMgr struct {
	mapDBLock sync.RWMutex
	mapDB map[string]*DBMongo
}

func ConstructorDBMongoMgr() *DBMongoMgr {
	dbCfg := server_common.GetAllDataBaseCfg()

	dbMgr := &DBMongoMgr{
		mapDB : make(map[string]*DBMongo),
	}

	for _, v := range dbCfg {
		db := ConstructorDBMongo(v.DataBaseUser, v.DataBasePWD, v.DataBaseIP, v.DataBasePort, v.DataBase)

		dbMgr.mapDBLock.Lock()
		dbMgr.mapDB[v.DataBaseAppName] = db
		dbMgr.mapDBLock.Unlock()
	}

	return dbMgr
}

func (pthis*DBMongoMgr)GetDBConnection(db string) *DBMongo {
	pthis.mapDBLock.RLock()
	defer pthis.mapDBLock.RUnlock()

	dbCon, ok := pthis.mapDB[db]
	if !ok || nil == dbCon{
		return nil
	}

	return dbCon
}
