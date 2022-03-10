package lin_common

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func goMemstatic() {
	for {
		//go vet ./…
		//GODEBUG=gctrace=1 ./xxx
		//GODEBUG=gctrace=1 ./xxx 2> gc.log
		runtime.GC()
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		fmt.Printf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes)\r\n", ms.Alloc, ms.HeapIdle, ms.HeapReleased)

		time.Sleep(time.Second * time.Duration(90))
	}
}

func ProfileInit() {
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil) //go tool pprof --text http://127.0.0.1:6060/debug/pprof/heap
	}()
	go goMemstatic()
}

func init(){
	//ProfileInit()
}