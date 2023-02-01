package common

// copy idea from https://github.com/panjf2000/ants
// https://github.com/fasthttp/

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type CALLBACK_FUNC_WORK func(CorPoolJobData)

type EN_CORPOOL_ERR int

const (
	EN_CORPOOL_ERR_no_free_worker = 1
	EN_CORPOOL_ERR_empty_lst      = 2
)

// coroutine pool err define
type CorPoolErr struct {
	str_     string
	errCode_ EN_CORPOOL_ERR
}

func (pthis *CorPoolErr) Error() string {
	return pthis.str_
}
func genCorpErr(ec EN_CORPOOL_ERR, param ...interface{}) error {
	err := &CorPoolErr{
		errCode_: ec,
		str_:     fmt.Sprint(param...),
	}
	return err
}

type EN_CORPOOL_JOBTYPE int

const (
	EN_CORPOOL_JOBTYPE_quit = 1
	EN_CORPOOL_JOBTYPE_user = 100
)

// coroutine pool job data
// JobType_ :
// JobData_ : the user job data
// JobCB_   : function call back in coroutine pool
type CorPoolJobData struct {
	JobType_ EN_CORPOOL_JOBTYPE
	JobData_ interface{}
	JobCB_   CALLBACK_FUNC_WORK
}

// coroutine pool worker define
type _corPoolWorker struct {
	corPool_  *CorPool
	jobChan_  chan CorPoolJobData
	workerID_ int64
	goID      uint64
}

type MAP_CORPOOLWORKER map[int64]*_corPoolWorker

// coroutine pool define
type CorPool struct {
	condPool_       *sync.Cond
	//condPoolTrigger bool

	WorkerFree_    list.List // *_corPoolWorker
	mapJobAll_       MAP_CORPOOLWORKER
	paramMaxCorCount int
	paramCheckCorCount int
	paramCheckInterval int
	corCount_      int
	wg_            sync.WaitGroup

	lastFreeCount_ int
}

func (worker *_corPoolWorker) _corWorkerDoJob(job *CorPoolJobData) {
	defer func() {
		err := recover()
		if err != nil {
			LogErr("recover get err:", err)
		}
	}()

	job.JobCB_(*job)
}

func (worker *_corPoolWorker) _go_CorWorker() {

	worker.goID = GetGID()

	worker.corPool_.wg_.Add(1)
	defer worker.corPool_.wg_.Done()

	COROUTINE_LOOP:
	for {
		jobData := <-worker.jobChan_
		if jobData.JobType_ == EN_CORPOOL_JOBTYPE_quit {
			worker._corWorkerDestroy()
			break COROUTINE_LOOP
		}

		// do job
		worker._corWorkerDoJob(&jobData)
		// add to cor pool free
		worker.corPool_.corPoolAddFreeWorker(worker)
	}
}
func (worker *_corPoolWorker) _corWorkerDestroy() {
	close(worker.jobChan_)
}

func (worker *_corPoolWorker) _corWorkerQuit() {
	worker.jobChan_ <- CorPoolJobData{
		JobType_: EN_CORPOOL_JOBTYPE_quit,
	}
}

func (pthis *CorPool) _go_CorPoolCheck() {
	for {
		pthis.condPool_.L.Lock()
		func(){
			defer func() {
				err := recover()
				if err != nil {
					LogErr("check cor pool err:", err)
				}
			}()
			curFreeCount := pthis.WorkerFree_.Len()
			if curFreeCount >= pthis.paramCheckCorCount && curFreeCount >= pthis.lastFreeCount_ && pthis.lastFreeCount_ >= pthis.paramCheckCorCount{
				LogDebug(" will quit some cor worker, curFreeCount:",
					curFreeCount, " lastFreeCount_:", pthis.lastFreeCount_, " paramCheckCorCount:", pthis.paramCheckCorCount)
				quitCount := curFreeCount / 2
				if quitCount < 1 {
					quitCount = 1
				}
				for i := 0; i < quitCount; i ++ {
					ele := pthis.WorkerFree_.Back()
					if ele == nil {
						break
					}
					pthis.WorkerFree_.Remove(ele)
					worker, ok := ele.Value.(*_corPoolWorker)
					if ok {
						worker._corWorkerQuit()
						delete(pthis.mapJobAll_, worker.workerID_)
					}
				}
			}
			pthis.lastFreeCount_ = pthis.WorkerFree_.Len()
		}()
		pthis.condPool_.L.Unlock()

		LogDebug(" check cor pool, curFreeCount:",
			pthis.WorkerFree_.Len(), " lastFreeCount_:", pthis.lastFreeCount_, " paramCheckCorCount:", pthis.paramCheckCorCount)
		time.Sleep(time.Second * time.Duration(pthis.paramCheckInterval))
	}
}


func (pthis *CorPool) corPoolAddFreeWorker(worker *_corPoolWorker) {

	pthis.condPool_.L.Lock()
	//pthis.condPoolTrigger = true
	bNeedSignal := false
	if pthis.corCount_ >= pthis.paramMaxCorCount && pthis.WorkerFree_.Len() == 0 {
		bNeedSignal = true
	}
	pthis.WorkerFree_.PushFront(worker)
	if bNeedSignal {
		pthis.condPool_.Broadcast()
	}
	pthis.condPool_.L.Unlock()
}


func (worker *_corPoolWorker) _corWorkerAddJob(job *CorPoolJobData) {
	worker.jobChan_ <- *job
}
// add a job to coroutine pool
func (pthis *CorPool) CorPoolAddJob(jobR *CorPoolJobData /* ready only */) error {
	pthis.condPool_.L.Lock()
	defer pthis.condPool_.L.Unlock()

/*	waitCount := 0
	tWaitBegin := time.Now().UnixMilli()*/

	for {
		if pthis.corCount_ >= pthis.paramMaxCorCount && pthis.WorkerFree_.Len() == 0 {
/*			lin_common.LogDebug("no worker, wait for free worker cor:",
				pthis.corCount_, " free:", pthis.WorkerFree_.Len())*/
			/*waitCount ++*/
			pthis.condPool_.Wait()
		} else {
			break
		}
	}

/*	tWaitEnd := time.Now().UnixMilli()
	if (tWaitEnd - tWaitBegin) > 50 * 1000 {
		lin_common.LogErr("wait too long:", tWaitEnd - tWaitBegin, " job data:", jobR.JobData_, " waitCount:", waitCount)
	}*/

	if pthis.corCount_ >= pthis.paramMaxCorCount && pthis.WorkerFree_.Len() == 0 {
		return genCorpErr(EN_CORPOOL_ERR_no_free_worker, "~~~~~~~~~no free work, cor:", pthis.corCount_, " free_len:", pthis.WorkerFree_.Len())
	}

	if pthis.WorkerFree_.Len() == 0 {
		newWorker := &_corPoolWorker{
			corPool_:  pthis,
			jobChan_:  make(chan CorPoolJobData, 100),
			workerID_: GenUUID64_V4(),
		}

		pthis.mapJobAll_[newWorker.workerID_] = newWorker
		pthis.corCount_++

		go newWorker._go_CorWorker()
		newWorker._corWorkerAddJob(jobR)
	} else {
		ele := pthis.WorkerFree_.Front() // put front, if the worker has worked once, it will work next time
		if ele != nil {
			pthis.WorkerFree_.Remove(ele)
			worker, ok := ele.Value.(*_corPoolWorker)
			if ok {
				worker._corWorkerAddJob(jobR)
			}
		} else {
			return genCorpErr(EN_CORPOOL_ERR_empty_lst, "list element is nil")
		}
	}

	return nil
}


// get a coroutine pool
func CorPoolInit(maxWorkerCount int, checkCount int, checkInterval int) *CorPool {
	if maxWorkerCount < 10 {
		maxWorkerCount = 10
	}
	if checkCount <=1 {
		checkCount = maxWorkerCount/2 + 1
	}
	cp := &CorPool{
		condPool_:      sync.NewCond(&sync.Mutex{}),
		paramMaxCorCount:   maxWorkerCount,
		paramCheckCorCount: checkCount,
		paramCheckInterval: checkInterval,
		corCount_:      0,
		mapJobAll_:     make(MAP_CORPOOLWORKER),
	}
	cp.WorkerFree_.Init()
	LogDebug("max worker count:", cp.paramMaxCorCount, " check count:", cp.paramCheckCorCount)

	go cp._go_CorPoolCheck()

	return cp
}

func (pthis *CorPool) corPoolUnitInter() {
	pthis.condPool_.L.Lock()
	defer pthis.condPool_.L.Unlock()

	// quit all worker
	for _, val := range pthis.mapJobAll_ {
		val._corWorkerQuit()
	}

	pthis.mapJobAll_ = make(MAP_CORPOOLWORKER)
}
func (pthis *CorPool) CorPoolUnit() {

	pthis.corPoolUnitInter()

	pthis.wg_.Wait()

	println("all worker quit")
}