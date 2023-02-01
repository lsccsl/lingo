package common

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestCorPoolInit(t *testing.T) {
	cp := CorPoolInit(10)

	cp.CorPoolAddJob(&CorPoolJobData{
		JobType_ : (EN_CORPOOL_JOBTYPE_user + 1),
		JobData_ : "interface{}",
		JobCB_   : func(jd CorPoolJobData){
			fmt.Println(jd.JobData_, jd.JobType_)
		},
	})

	cp.CorPoolAddJob(&CorPoolJobData{
		JobType_ : (EN_CORPOOL_JOBTYPE_user + 2),
		JobData_ : "aaaaa{}",
		JobCB_   : func(jd CorPoolJobData){
			fmt.Println(jd.JobData_, jd.JobType_)
		},
	})

	cp.CorPoolAddJob(&CorPoolJobData{
		JobType_ : (EN_CORPOOL_JOBTYPE_user + 3),
		JobData_ : "aaaaabbb{}",
		JobCB_   : func(jd CorPoolJobData){
			fmt.Println(jd.JobData_, jd.JobType_)
		},
	})

	var g_count int64 = 0
	tbegin := time.Now().Unix()
	for i := 0; i < 10 * 1000 * 1000; i ++ {
		str := fmt.Sprint("aaaaa", i)
		cp.CorPoolAddJob(&CorPoolJobData{
			JobType_ : (EN_CORPOOL_JOBTYPE_user + EN_CORPOOL_JOBTYPE(i)),
			JobData_ : str,
			JobCB_   : func(jd CorPoolJobData){
				atomic.AddInt64(&g_count, 1)
				if atomic.LoadInt64(&g_count) % 100000 == 0 {
					fmt.Println(jd.JobData_, g_count)
				}
			},
		})
	}

	cp.CorPoolUnit()
	tend := time.Now().Unix()
	fmt.Println("end", tend - tbegin)
}

