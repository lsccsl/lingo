package main

import (
	"fmt"
	. "lin/lin_common"
	"sync"
	"sync/atomic"
	"time"
)

func test_chan() {
	fmt.Println("begine test chang")
	c := make(chan int, 100)

	go func() {
		c <- 1
		c <- 2
		c <- 0
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {

		for v := range c {
			fmt.Println(v)
			if v == 0 {
				break
			}
		}
		fmt.Println("end range coroutine")
		wg.Done()
	}()

	wg.Wait()

	fmt.Println("end test chang")
}

func test_cor_pool(max_count int, work_sleep int, loop_count int) {
	cp := CorPoolInit(max_count)

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
	for i := 0; i < loop_count/*10 * 1000 * 1000*/; i ++ {
		str := fmt.Sprint("aaaaa", i)
		cp.CorPoolAddJob(&CorPoolJobData{
			JobType_ : (EN_CORPOOL_JOBTYPE_user + EN_CORPOOL_JOBTYPE(i)),
			JobData_ : str,
			JobCB_   : func(jd CorPoolJobData){
				atomic.AddInt64(&g_count, 1)
				if atomic.LoadInt64(&g_count) % 1000000 == 0 {
					fmt.Println(jd.JobData_, g_count)
				}
				if work_sleep > 0 {
					time.Sleep(time.Duration(work_sleep) * time.Microsecond)
				}
			},
		})
	}

	cp.CorPoolUnit()
	tend := time.Now().Unix()
	fmt.Println("end", tend - tbegin)

	time.Sleep(3 * time.Second)
}

func main() {
	//test_chan()
	test_cor_pool(1000, 0, 10 * 1000 * 1000)
	test_cor_pool(10, 10, 1000 * 1000)
}