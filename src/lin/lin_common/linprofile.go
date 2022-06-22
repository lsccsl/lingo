package lin_common

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"time"
)

func goStatic() {
	for {
		//go vet ./…
		//GODEBUG=gctrace=1 ./xxx
		//GODEBUG=gctrace=1 ./xxx 2> gc.log
		// cat /var/log/messages
		runtime.GC()
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		//LogDebug("Alloc:", ms.Alloc, "(bytes) HeapIdle:", ms.HeapIdle, "(bytes) HeapReleased:", ms.HeapReleased, "(bytes)", " coroutine:", runtime.NumGoroutine())

		time.Sleep(time.Second * time.Duration(180))
	}
}

func ProfileInit(needStatic bool, profilePort int) {
	runtime.GOMAXPROCS(16)
	LogDebug(runtime.NumCPU())
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	go func() {
		addr := "0.0.0.0:" + strconv.Itoa(profilePort)
		http.ListenAndServe(addr, nil) //go tool pprof --text http://127.0.0.1:6060/debug/pprof/heap
	}()
	if needStatic {
		go goStatic()
	}
}

func init(){
	//ProfileInit()
}