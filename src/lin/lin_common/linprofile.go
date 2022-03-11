package lin_common

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func goStatic() {
	for {
		//go vet ./…
		//GODEBUG=gctrace=1 ./xxx
		//GODEBUG=gctrace=1 ./xxx 2> gc.log
		runtime.GC()
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		LogDebug("Alloc:", ms.Alloc, "(bytes) HeapIdle:", ms.HeapIdle, "(bytes) HeapReleased:", ms.HeapReleased, "(bytes)", " coroutine:", runtime.NumGoroutine())

		time.Sleep(time.Second * time.Duration(10))
	}
}

func ProfileInit() {
	//runtime.GOMAXPROCS(8)
	LogDebug(runtime.NumCPU())
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil) //go tool pprof --text http://127.0.0.1:6060/debug/pprof/heap
	}()
	go goStatic()
}

func init(){
	//ProfileInit()
}