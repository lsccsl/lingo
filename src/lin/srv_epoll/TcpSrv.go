package main

import (
	"lin/lin_common"
	"time"
)

type TcpSrv struct {
	fd lin_common.FD_DEF

	timerConnClose * time.Timer
	durationClose time.Duration
	pu *eSrvMgrProcessUnit
}


