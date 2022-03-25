package cor_pool

// copy idea from https://github.com/panjf2000/ants
// https://github.com/fasthttp/

import (
	"container/list"
	"fmt"
	"lin/lin_common"
	"sync"
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
	workerID_ int
	goID      uint64
}

// coroutine pool define
type CorPool struct {
	lockPool_ sync.Mutex

	condPool_       *sync.Cond
	condPoolTrigger bool

	WorkerFree_    list.List // *_corPoolWorker
	mapJobAll_     map[int]*_corPoolWorker
	maxCorCount_   int
	checkCorCount_ int
	corCount_      int
	wg_            sync.WaitGroup
}

func (worker *_corPoolWorker) _corWorkerDoJob(job *CorPoolJobData) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr("recover get err:", err)
		}
	}()

	job.JobCB_(*job)
}

func (worker *_corPoolWorker) _go_CorWorker() {

	worker.goID = lin_common.GetGID()

	worker.corPool_.wg_.Add(1)
	defer worker.corPool_.wg_.Done()

	COROUTINE_LOOP:
	for {
		//println("chan msg count", len(worker.jobChan_))
		jobData := <-worker.jobChan_
		//println("get job", jobData.JobType_, "worker:", worker.workerID_, lin_common.GetGID())
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



func (pthis *CorPool) corPoolAddFreeWorker(worker *_corPoolWorker) {
	pthis.lockPool_.Lock()
	defer pthis.lockPool_.Unlock()

	bNeedTriggerSignal := false
	if pthis.corCount_ >= pthis.maxCorCount_ && pthis.WorkerFree_.Len() == 0 {
		bNeedTriggerSignal = true
	}
	pthis.WorkerFree_.PushFront(worker)

	if (bNeedTriggerSignal) {
		//lin_common.LogDebug("trigger signal")
		pthis.condPool_.L.Lock()
		pthis.condPoolTrigger = true
		pthis.condPool_.Signal()
		pthis.condPool_.L.Unlock()
	}
}


func (worker *_corPoolWorker) _corWorkerAddJob(job *CorPoolJobData) {
	worker.jobChan_ <- *job
}
// add a job to coroutine pool
func (pthis *CorPool) CorPoolAddJob(jobR *CorPoolJobData /* ready only */) error {
	{
		pthis.lockPool_.Lock()

		if pthis.corCount_ >= pthis.maxCorCount_ && pthis.WorkerFree_.Len() == 0 {
			pthis.lockPool_.Unlock()
			lin_common.LogDebug("no worker, wait for free worker ~~~~~~~~~~~~~~~~~~~ cor:", pthis.corCount_, " free:", pthis.WorkerFree_.Len())

			pthis.condPool_.L.Lock()
			if !pthis.condPoolTrigger {
				//println("wait signal")
				pthis.condPool_.Wait()
			}
			pthis.condPoolTrigger = false
			pthis.condPool_.L.Unlock()

			pthis.lockPool_.Lock()
		}

		defer func() {
			pthis.lockPool_.Unlock()
			//println("unlock")
		}()

		if pthis.corCount_ >= pthis.maxCorCount_ && pthis.WorkerFree_.Len() == 0 {
			return genCorpErr(EN_CORPOOL_ERR_no_free_worker, "no free work, cor:", pthis.corCount_, " free_len:", pthis.WorkerFree_.Len())
		}

		if pthis.WorkerFree_.Len() == 0 {

			newWorker := &_corPoolWorker{
				corPool_: pthis,
				jobChan_: make(chan CorPoolJobData, 100),
			}

			pthis.mapJobAll_[pthis.corCount_] = newWorker
			newWorker.workerID_ = pthis.corCount_
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
	}

	return nil
}


// get a coroutine pool
func CorPoolInit(maxWorkerCount int) *CorPool {
	cp := &CorPool{
		condPool_:      sync.NewCond(&sync.Mutex{}),
		maxCorCount_:   maxWorkerCount,
		checkCorCount_: (maxWorkerCount/2 + 1),
		corCount_:      0,
		mapJobAll_:     make(map[int]*_corPoolWorker),
	}
	cp.WorkerFree_.Init()
	lin_common.LogDebug("max worker count:", maxWorkerCount)

	return cp
}

func (pthis *CorPool) corPoolUnitInter() {
	pthis.lockPool_.Lock()
	defer pthis.lockPool_.Unlock()

	// quit all worker
	for _, val := range pthis.mapJobAll_ {
		val._corWorkerQuit()
	}

	pthis.mapJobAll_ = make(map[int]*_corPoolWorker)
}
func (pthis *CorPool) CorPoolUnit() {

	pthis.corPoolUnitInter()

	pthis.wg_.Wait()

	println("all worker quit")
}