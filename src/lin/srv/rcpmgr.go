package main

import "sync"

type RPCReq struct {
	rpcID int64
	chNtf chan interface{}
}
type MAP_RPC_REQ map[int64/* rpc msg id */]*RPCReq
type RPCManager struct {
	mapRPCReqMutex sync.Mutex
	mapRPCReq MAP_RPC_REQ
}

func ConstructRPCManager() *RPCManager {
	rmgr := &RPCManager{
		mapRPCReq:make(MAP_RPC_REQ),
	}

	return rmgr
}

func (pthis*RPCManager)RPCManagerAddReq(rpcID int64) * RPCReq {
	pthis.mapRPCReqMutex.Lock()
	defer pthis.mapRPCReqMutex.Unlock()

	val, _ := pthis.mapRPCReq[rpcID]
	if val != nil {
		return nil
	}
	rreq := &RPCReq{rpcID: rpcID,
		chNtf: make(chan interface{}, 10),
	}
	pthis.mapRPCReq[rpcID] = rreq
	return rreq
}

func (pthis*RPCManager)RPCManagerFindReq(rpcID int64) * RPCReq {
	pthis.mapRPCReqMutex.Lock()
	defer pthis.mapRPCReqMutex.Unlock()
	val, _ := pthis.mapRPCReq[rpcID]
	return val
}

func (pthis*RPCManager)RPCManagerDelReq(rpcID int64) {
	pthis.mapRPCReqMutex.Lock()
	defer pthis.mapRPCReqMutex.Unlock()
	val, _ := pthis.mapRPCReq[rpcID]
	if val != nil {
		if val.chNtf != nil {
			close(val.chNtf)
		}
	}
	delete(pthis.mapRPCReq, rpcID)
}
