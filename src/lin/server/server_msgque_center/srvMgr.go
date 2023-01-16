package main

import "lin/lin_common"

type srvMgr struct {

}

type srvOne struct {
	srvID int64
	fdSrv lin_common.FD_DEF
	ip string
	port int32
}

func ConstructSrvMgr() *srvMgr {
	return &srvMgr{}
}
